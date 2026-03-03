package common

import (
	"bytes"
	"context"
	"fmt"
	"m4s-converter/conver"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func (c *Config) overlay() string {
	if c.Overlay {
		return "-y"
	}
	return "-n"
}

func (c *Config) Composition(videoFile, audioFile, outputFile string) error {
	var args []string
	if c.Overlay {
		args = append(args, "-force")
	}

	args = append(args, "-charset", "utf8")

	tags := fmt.Sprintf("title=%s:artist=%s:album=%s", c.GroupId, c.Uid, c.ItemId)
	args = append(args, "-tags", tags)
	args = append(args,
		"-cprt", c.ItemId,
		"-add", videoFile+"#video",
		"-add", audioFile+"#audio",
		"-new", outputFile)

	cmd := exec.Command(c.GPACPath, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	if c.AssPath != "" {
		assFile := strings.ReplaceAll(outputFile, conver.Mp4Suffix, conver.AssSuffix)
		_ = c.copyFile(c.AssPath, assFile)
	}

	if err := cmd.Run(); err != nil {
		logrus.Errorf("合成视频文件失败:%s\n%s", outputFile, stdout.String())
		return err
	}

	logrus.Info("已合成视频文件:", outputFile)
	return nil
}

func (c *Config) FindM4sFiles(src string, info os.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if !strings.HasSuffix(info.Name(), conver.M4sSuffix) {
		return nil
	}

	var dst string
	videoId, audioId := GetVAId(src)
	if videoId != "" && audioId != "" {
		if strings.Contains(info.Name(), audioId) {
			dst = strings.ReplaceAll(src, conver.M4sSuffix, conver.AudioSuffix)
		} else {
			dst = strings.ReplaceAll(src, conver.M4sSuffix, conver.VideoSuffix)
		}
	}

	if err = c.M4sToAV(src, dst); err != nil {
		MessageBox(fmt.Sprintf("%v 转换异常：%v", src, err))
		return err
	}
	logrus.Info("已将m4s转换为音视频文件: ", strings.TrimPrefix(dst, c.CachePath))
	return nil
}

func (c *Config) GetAudioAndVideo(ctx context.Context, cachePath string) (string, string, error) {
	var video, audio string

	entries, err := os.ReadDir(cachePath)
	if err != nil {
		return "", "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			childVideo, childAudio, err := c.GetAudioAndVideo(ctx, filepath.Join(cachePath, entry.Name()))
			if err == nil && childVideo != "" && childAudio != "" {
				video = childVideo
				audio = childAudio
				break
			}
			continue
		}

		fileName := entry.Name()
		if strings.HasSuffix(fileName, conver.VideoSuffix) {
			video = filepath.Join(cachePath, fileName)
		}
		if strings.HasSuffix(fileName, conver.AudioSuffix) {
			audio = filepath.Join(cachePath, fileName)
		}
	}

	if video == "" || audio == "" {
		return "", "", fmt.Errorf("找不到音频或视频文件: %s", cachePath)
	}

	if !c.AssOFF {
		oldVideo := c.video
		c.video = video
		c.downloadXml(ctx)
		c.video = oldVideo
	}

	return video, audio, nil
}

func joinUrl(cid string) string {
	return "https://comment.bilibili.com/" + cid + conver.XmlSuffix
}

func joinXmlUrl(cid string) string {
	return "https://api.bilibili.com/x/v1/dm/list.so?oid=" + cid
}

func (c *Config) downloadXml(ctx context.Context) {
	dirPath := filepath.Dir(c.video)
	dirName := filepath.Base(dirPath)

	if len(dirName) < 6 {
		danmakuXml := filepath.Join(filepath.Dir(dirPath), conver.DanmakuXml)
		if Size(danmakuXml) != 0 {
			c.AssPath = conver.Xml2Ass(danmakuXml)
		}
		return
	}

	xmlPath := filepath.Join(dirPath, dirName+conver.XmlSuffix)
	if Size(xmlPath) != 0 {
		c.AssPath = conver.Xml2Ass(xmlPath)
		return
	}

	if err := downloadFile(ctx, joinUrl(dirName), xmlPath); err != nil {
		if err := downloadFile(ctx, joinXmlUrl(dirName), xmlPath); err != nil {
			logrus.Warn("弹幕文件下载失败:", joinUrl(dirName))
			return
		}
	}
	c.AssPath = conver.Xml2Ass(xmlPath)
}

func getVAIdFromAndroidMode(patch string) (string, string) {
	if filepath.Base(filepath.Dir(patch)) != "80" {
		logrus.Warnln("找不到.playurl文件,切换到Android模式解析entry.json文件")
	}

	androidPEJ := filepath.Join(filepath.Dir(filepath.Dir(patch)), conver.PlayEntryJson)
	puDate, err := os.ReadFile(androidPEJ)
	if err != nil {
		logrus.Error("找不到entry.json文件!")
		return "", ""
	}

	status := gjson.GetBytes(puDate, "page_data.download_title").String()
	if status != "completed" && status != "视频已缓存完成" && status != "" {
		logrus.Error("跳过未缓存完成的视频", status)
		return "", ""
	}

	return "video.m4s", "audio.m4s"
}

func extractIDsFromPlayUrl(puByte []byte) (string, string) {
	var p gjson.Result
	if p = gjson.GetBytes(puByte, "data"); !p.Exists() {
		p = gjson.GetBytes(puByte, "result")
	}
	if p.Exists() {
		return p.Get("dash.video|@reverse|0.id").String(), p.Get("dash.audio|@reverse|0.id").String()
	}
	return "", ""
}

func (c *Config) PanicHandler() {
	if e := recover(); e != nil {
		fmt.Print("按回车键退出...")
		_, _ = fmt.Scanln()
	}
}

func (c *Config) isIdenticalFileExists(dirPath string, videoPath string, audioPath string) (bool, string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		logrus.Errorf("读取目录失败: %v", err)
		return false, ""
	}

	inputHash := c.calculateCombinedHash(videoPath, audioPath)
	if inputHash == "" {
		return checkByFileSize(files, dirPath, videoPath, audioPath)
	}

	return checkByHash(files, dirPath, inputHash, c)
}

func checkByFileSize(files []os.DirEntry, dirPath, videoPath, audioPath string) (bool, string) {
	videoInfo, err := os.Stat(videoPath)
	if err != nil {
		return false, ""
	}

	audioInfo, err := os.Stat(audioPath)
	if err != nil {
		return false, ""
	}

	expectedSize := videoInfo.Size() + audioInfo.Size()

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".mp4") {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		if abs(int64(fileInfo.Size())-expectedSize) <= sizeTolerance {
			return true, filePath
		}
	}

	return false, ""
}

func checkByHash(files []os.DirEntry, dirPath, inputHash string, c *Config) (bool, string) {
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".mp4") {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		hashFilePath := strings.ReplaceAll(filePath, ".mp4", ".hash")
		if hashContent, err := os.ReadFile(hashFilePath); err == nil && string(hashContent) == inputHash {
			return true, filePath
		}

		if metadata, err := c.getMp4Metadata(filePath); err == nil {
			if metadata["title"] == c.GroupId && metadata["artist"] == c.Uid && metadata["album"] == c.ItemId {
				return true, filePath
			}
		}
	}

	return false, ""
}
