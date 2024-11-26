package common

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"m4s-converter/conver"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (c *Config) Synthesis() {
	begin := time.Now().Unix()
	logrus.Println("查找缓存目录下可转换的文件...")
	// 查找m4s文件，并转换为mp4和mp3
	if err := filepath.WalkDir(c.CachePath, c.FindM4sFiles); err != nil {
		MessageBox(fmt.Sprintf("查找并转换 m4s 文件异常：%v", err))
		c.wait()
	}

	dirs, err := GetCacheDir(c.CachePath) // 缓存根目录模式
	if err != nil {
		MessageBox(fmt.Sprintf("找不到 BiliBili 的缓存目录：%v", err))
		c.wait()
	}

	if dirs == nil {
		// 判断非缓存根目录时，验证是否为子目录
		if Exist(filepath.Join(c.CachePath, conver.VideoInfoSuffix)) ||
			Exist(filepath.Join(c.CachePath, conver.VideoInfoJson)) {
			dirs = append(dirs, c.CachePath)
		}
	}

	// 合成音视频文件
	c.OutputDir = filepath.Join(c.CachePath, "output")
	var outputFiles []string
	var skipFilePaths []string
	for _, v := range dirs {
		video, audio, e := c.GetAudioAndVideo(v)
		if e != nil {
			logrus.Error("找不到已修复的音频和视频文件:", err)
			continue
		}
		info := filepath.Join(v, conver.VideoInfoJson)
		if !Exist(info) {
			info = filepath.Join(v, conver.VideoInfoSuffix)
			if !Exist(info) {
				info = filepath.Join(v, conver.PlayEntryJson)
				if !Exist(info) {
					continue
				}
			}
		}
		infoStr, e := os.ReadFile(info)
		if e != nil {
			logrus.Error("找不到包含视频信息的info相关文件: ", info)
			continue
		}
		js, e := simplejson.NewJson(infoStr)
		if e != nil {
			logrus.Error("videoInfo相关文件解析失败: ", info)
			continue
		}

		groupTitle := Filter(js.Get("groupTitle").String())
		groupTitle = null2Str(groupTitle, Filter(js.Get("owner_name").String()))

		title := Filter(js.Get("page_data").Get("download_subtitle").String())
		title = null2Str(title, Filter(js.Get("title").String()))

		uname := Filter(js.Get("uname").String())
		uname = null2Str(uname, Filter(js.Get("title").String()))

		status := Filter(js.Get("status").String())
		status = null2Str(status, Filter(js.Get("page_data").Get("download_title").String()))

		itemId, e := js.Get("itemId").Int()
		if itemId == 0 || e != nil {
			itemId, _ = js.Get("owner_id").Int()
		}
		c.ItemId = strconv.Itoa(itemId)

		if status != "completed" && status != "视频已缓存完成" {
			skipFilePaths = append(skipFilePaths, v)
			logrus.Warn("未缓存完成,跳过合成", v, title+"-"+uname)
			continue
		}
		if !Exist(c.OutputDir) {
			_ = os.MkdirAll(c.OutputDir, os.ModePerm)
		}
		groupPath := groupTitle + "-" + uname
		groupDir := filepath.Join(c.OutputDir, groupPath)
		if !Exist(groupDir) {
			if err = os.MkdirAll(groupDir, os.ModePerm); err != nil {
				MessageBox("无法创建目录：" + groupDir)
				c.wait()
			}
		}
		mp4Name := title + conver.Mp4Suffix
		outputFile := filepath.Join(groupDir, mp4Name)
		if c.Skip || Exist(outputFile) && c.findMp4Info(outputFile, c.ItemId) {
			logrus.Warn("跳过合成完全相同的视频:", filepath.Join(groupPath, mp4Name))
			continue
		}
		if Exist(outputFile) && !c.Overlay {
			mp4Name = title + c.ItemId + conver.Mp4Suffix
			outputFile = filepath.Join(groupDir, mp4Name)
		}
		if c.findMp4Info(outputFile, c.ItemId) {
			logrus.Warn("跳过合成完全相同的视频:", filepath.Join(groupPath, mp4Name))
			continue
		}

		if er := c.Composition(video, audio, outputFile); er != nil {
			logrus.Errorf("%s 合成失败", filepath.Base(outputFile))
			continue
		}
		outputFiles = append(outputFiles, filepath.Join(groupPath, mp4Name))
	}

	end := time.Now().Unix()
	logrus.Print("==========================================")
	if skipFilePaths != nil {
		logrus.Print("跳过的目录:\n" + strings.Join(skipFilePaths, "\n"))
	}
	if outputFiles != nil {
		logrus.Printf("# 输出目录:\n%s", color.CyanString(c.OutputDir))
		logrus.Printf("# 合成的文件:\n%s", color.CyanString(strings.Join(outputFiles, "\n")))
		// 打开合成文件目录
		go OpenFolder(c.OutputDir)
	} else {
		logrus.Warn("未合成任何文件！")
	}
	logrus.Print("==========================================")
	logrus.Print("已完成合成任务，耗时:", end-begin, "秒")
	c.wait()
}

func (c *Config) findMp4Info(fp, sub string) bool {
	if !Exist(c.GPACPath) {
		return false
	}
	ret, err := exec.Command(c.GPACPath, "-info", fp).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(ret), sub)
}

func null2Str(s string, value string) string {
	if s != "" {
		return s
	}
	return value
}

func (c *Config) wait() {
	fmt.Print("按[回车]键退出...")
	_, _ = fmt.Scanln()
	os.Exit(0)
}
