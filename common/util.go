package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"m4s-converter/conver"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	utils "github.com/mzky/utils/common"
	"github.com/ncruces/zenity"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Config struct {
	FFMpegPath string
	CachePath  string
	Overlay    bool
	Skip       bool
	AssPath    string
	AssOFF     bool
	OutputDir  string
	GPACPath   string
	video      string
	audio      string
	ItemId     string
}

func (c *Config) overlay() string {
	if c.Overlay {
		return "-y"
	}
	return "-n"
}
func (c *Config) Composition(videoFile, audioFile, outputFile string) error {
	var cmd *exec.Cmd
	if c.GPACPath != "" {
		cmd = exec.Command(c.GPACPath,
			// "-quiet", // 仅打印异常日志
			"-cprt", c.ItemId,
			"-add", videoFile+"#video",
			"-add", audioFile+"#audio",
			"-new", outputFile)
	} else {
		// 构建FFmpeg命令行参数
		var args []string
		args = append(args,
			"-i", videoFile,
			"-i", audioFile,
			"-c:v", "copy", // video不指定编解码，使用 BiliBili 原有编码
			"-c:a", "copy", // audio不指定编解码可能会导致音视频不同步
			// "-strict", "experimental", // 宽松编码控制器
			"-vsync", "2", // 根据音频流调整视频帧率
			// "-async", "1", // 强制音频流与视频流同步，通过丢弃或重复音频样本以匹配视频流
			"-shortest",   // 输出文件的长度与较短的那个流相同，防止过长的流导致不同步
			"-map", "0:v", // 指定从第一个输入文件中选择视频流
			"-map", "1:a", // 从第二个输入文件中选择音频流
			c.overlay(),               // 是否覆盖已存在视频
			"-movflags", "+faststart", // 启用faststart可以让视频在网络传输时更快地开始播放
			outputFile,
			"-hide_banner", // 隐藏版本信息和版权声明
			"-stats",       // 只显示统计信息
		)
		cmd = exec.Command(c.FFMpegPath, args...)
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	if c.AssPath != "" {
		assFile := strings.ReplaceAll(outputFile, conver.Mp4Suffix, conver.AssSuffix)
		_ = c.copyFile(c.AssPath, assFile)
	}

	// 等待命令执行完成
	if err := cmd.Run(); err != nil {
		logrus.Errorf("合成视频文件失败:%s\n%s", filepath.Base(outputFile), stdout.String())
		return nil
	}

	logrus.Info("已合成视频文件:", filepath.Base(outputFile))
	return nil
}

func (c *Config) FindM4sFiles(src string, info os.DirEntry, err error) error {
	if err != nil {
		return err
	}
	// 查找.m4s文件
	if strings.HasSuffix(info.Name(), conver.M4sSuffix) {
		var dst string
		videoId, audioId := GetVAId(src)
		if videoId != "" && audioId != "" {
			if strings.Contains(info.Name(), audioId) { // 音频文件
				dst = strings.ReplaceAll(src, conver.M4sSuffix, conver.AudioSuffix)
			} else {
				dst = strings.ReplaceAll(src, conver.M4sSuffix, conver.VideoSuffix)
			}
		}

		if err = c.M4sToAV(src, dst); err != nil {
			MessageBox(fmt.Sprintf("%v 转换异常：%v", src, err))
			return err
		}
		logrus.Info("已将m4s转换为音视频文件: ", strings.TrimLeft(dst, c.CachePath))
	}
	return nil
}

func GetCacheDir(cachePath string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != cachePath {
			if !strings.Contains(path, "output") {
				dirs = append(dirs, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirs, nil
}

func joinUrl(cid string) string {
	return "https://comment.bilibili.com/" + cid + conver.XmlSuffix
}
func joinXmlUrl(cid string) string {
	return "https://api.bilibili.com/x/v1/dm/list.so?oid=" + cid
}

// GetAudioAndVideo 从给定的缓存路径中查找音频和视频文件，并尝试下载并转换xml弹幕为ass格式
// 参数:
// - cachePath: 缓存路径，用于搜索音频、视频文件以及存储下载的弹幕文件
// 返回值:
// - video: 查找到的视频文件路径
// - audio: 查找到的音频文件路径
// - error: 在搜索、下载或转换过程中遇到的任何错误
func (c *Config) GetAudioAndVideo(cachePath string) (string, string, error) {
	// 遍历给定路径下的所有文件和目录
	if err := filepath.Walk(cachePath, c.findAV); err != nil {
		return "", "", err // 如果遍历过程中发生错误，返回错误信息
	}
	// 下载弹幕文件
	if !c.AssOFF {
		c.downloadXml()
	}
	return c.video, c.audio, nil // 返回找到的视频和音频文件路径
}

func (c *Config) findAV(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err // 如果遇到错误，立即返回
	}
	if !info.IsDir() {
		// 如果是文件，检查是否为视频或音频文件
		if strings.Contains(path, conver.VideoSuffix) {
			c.video = path // 找到视频文件
		}
		if strings.Contains(path, conver.AudioSuffix) {
			c.audio = path // 找到音频文件
		}
	}
	return nil
}
func (c *Config) copyFile(src, dst string) error {
	// 打开源文件
	srcFile, e := os.Open(src)
	if e != nil {
		logrus.Errorf("打开源文件失败: %v", e)
		return e
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, e := os.Create(dst)
	if e != nil {
		logrus.Errorf("创建目标文件失败: %v", e)
		return e
	}
	defer dstFile.Close()

	// 读取前 9 个字节
	data := make([]byte, 9)
	if _, err := io.ReadAtLeast(srcFile, data, 9); err != nil {
		logrus.Errorf("读取文件头失败: %v", err)
		return e
	}

	// 检查前 9 个字节是否为 '0'
	if string(data) != "000000000" {
		// 如果前 9 个字节不为 '0'，写入这些字节
		_, _ = dstFile.Write(data)
	}

	// 使用缓冲读取器逐块读取并写入文件
	if _, err := io.Copy(bufio.NewWriter(dstFile), bufio.NewReader(srcFile)); err != nil {
		logrus.Errorf("读取或写入文件失败: %v", err)
		return err
	}
	return nil
}

func (c *Config) M4sToAV(src, dst string) error {
	return c.copyFile(src, dst)
}

// GetCachePath 获取用户视频缓存路径
func (c *Config) GetCachePath() {
	if c.findM4sFiles() != nil {
		MessageBox("BiliBili缓存路径 " + c.CachePath + " 未找到缓存文件, \n请重新选择 BiliBili 缓存文件路径！")
		c.SelectDirectory()
		return
	}
	logrus.Info("选择的 BiliBili 缓存目录为: ", c.CachePath)
	return
}

func Size(path string) int64 {
	if utils.IsExist(path) {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return 0
		}
		return fileInfo.Size()
	}
	return 0
}

// Filter 过滤文件名
func Filter(name string, err error) string {
	if err != nil || name == "" {
		return ""
	}
	name = strings.ReplaceAll(name, "（", "(")
	name = strings.ReplaceAll(name, "）", ")")
	name = strings.ReplaceAll(name, "<", "《")
	name = strings.ReplaceAll(name, ">", "》")
	name = strings.ReplaceAll(name, `\`, "#")
	name = strings.ReplaceAll(name, `"`, "'")
	name = strings.ReplaceAll(name, "/", "#")
	name = strings.ReplaceAll(name, "|", "_")
	name = strings.ReplaceAll(name, "?", "？")
	name = strings.ReplaceAll(name, "*", "-")
	name = strings.ReplaceAll(name, "【", "[")
	name = strings.ReplaceAll(name, "】", "]")
	name = strings.ReplaceAll(name, ":", "：")
	name = strings.ReplaceAll(name, " ", "")

	return strings.TrimSpace(name)
}

func (c *Config) PanicHandler() {
	if e := recover(); e != nil {
		fmt.Print("按回车键退出...")
		_, _ = fmt.Scanln()
	}
}

func MessageBox(text string) {
	_ = zenity.Warning(text, zenity.Title("提示"), zenity.Width(400))
}

// findM4sFiles 检查目录及其子目录下是否存在m4s文件
func (c *Config) findM4sFiles() error {
	var m4sFiles []string
	err := filepath.Walk(c.CachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Warnf("查找bilibili缓存目录异常: %s", path)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == conver.M4sSuffix {
			m4sFiles = append(m4sFiles, path)
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(m4sFiles) == 0 {
		return fmt.Errorf("缓存目录找不到m4s文件: %s", c.CachePath)
	}
	return nil
}

// SelectDirectory 选择 BiliBili 缓存目录
func (c *Config) SelectDirectory() {
	var err error
	c.CachePath, err = zenity.SelectFile(zenity.Title("请选择 BiliBili 缓存目录"), zenity.Directory())
	if c.CachePath == "" || err != nil {
		logrus.Warn("关闭对话框后自动退出程序")
		os.Exit(1)
	}

	if c.findM4sFiles() == nil {
		logrus.Info("选择的 BiliBili 缓存目录为:", c.CachePath)
		return
	}
	MessageBox("选择的 BiliBili 缓存目录内找不到m4s文件，请重新选择！")
	c.SelectDirectory()
}

// SelectGPACPath 选择 GPACPath文件
func (c *Config) SelectGPACPath() {
	var err error
	c.GPACPath, err = zenity.SelectFile(zenity.Title("请选择 GPAC 的 mp4box 文件"))
	if c.GPACPath == "" || err != nil {
		logrus.Warn("关闭对话框后自动退出程序")
		os.Exit(1)
	}

	if utils.IsExist(c.GPACPath) {
		logrus.Info("选择 GPAC 的 mp4box 文件为:", c.CachePath)
		return
	}
	MessageBox("选择 GPAC 的 mp4box 文件不存在，请重新选择！")
	c.SelectGPACPath()
}

func (c *Config) SelectFFMpegPath() {
	var err error
	c.FFMpegPath, err = zenity.SelectFile(zenity.Title("请选择 FFMpeg 文件"))
	if c.FFMpegPath == "" || err != nil {
		logrus.Warn("关闭对话框后自动退出程序")
		os.Exit(1)
	}

	if utils.IsExist(c.FFMpegPath) {
		logrus.Info("选择 FFMpeg 文件为:", c.CachePath)
		return
	}
	MessageBox("选择 FFMpeg 文件不存在，请重新选择！")
	c.SelectFFMpegPath()
}

// 如果是目录，尝试下载并转换xml弹幕为ass格式
func (c *Config) downloadXml() {
	dirPath := filepath.Dir(c.video)
	dirName := filepath.Base(dirPath)

	if len(dirName) < 6 { // Android嵌套目录，音视频目录为80
		danmakuXml := filepath.Join(filepath.Dir(dirPath), conver.DanmakuXml)
		if Size(danmakuXml) != 0 {
			c.AssPath = conver.Xml2Ass(danmakuXml) // 转换xml弹幕文件为ass格式
		}
		return
	}
	xmlPath := filepath.Join(dirPath, dirName+conver.XmlSuffix)
	if Size(xmlPath) != 0 {
		c.AssPath = conver.Xml2Ass(xmlPath) // 转换xml弹幕文件为ass格式
		return
	}
	if e := downloadFile(joinUrl(dirName), xmlPath); e != nil {
		if downloadFile(joinXmlUrl(dirName), xmlPath) != nil {
			logrus.Warn("弹幕文件下载失败:", joinUrl(dirName))
			return
		}
	}
	c.AssPath = conver.Xml2Ass(xmlPath) // 转换xml弹幕文件为ass格式
}

// GetVAId 返回.playurl文件中视频文件或音频文件件数组
func GetVAId(patch string) (videoID string, audioID string) {
	pu := filepath.Join(filepath.Dir(patch), conver.PlayUrlSuffix)
	puByte, e := os.ReadFile(pu)
	if e == nil {
		/*
			视频：
			data.dash.video[0].id
			data.dash.audio[0].id
			番剧：
			result.dash.video[0].id  80  需要加上30000，实际30080.m4s
			result.dash.audio[0].id  30280
		*/
		var p gjson.Result
		if p = gjson.GetBytes(puByte, "data"); !p.Exists() {
			p = gjson.GetBytes(puByte, "result")
		}
		if p.Exists() {
			return p.Get("dash.video|@reverse|0.id").String(), p.Get("dash.audio|@reverse|0.id").String()
		}
		return "", ""
	}
	if filepath.Base(filepath.Dir(patch)) != "80" {
		logrus.Warnln("找不到.playurl文件,切换到Android模式解析entry.json文件")
	}
	androidPEJ := filepath.Join(filepath.Dir(filepath.Dir(patch)), conver.PlayEntryJson)
	puDate, e := os.ReadFile(androidPEJ)
	if e != nil {
		logrus.Error("找不到entry.json文件!")
		return
	}
	status := gjson.GetBytes(puDate, "page_data.download_title").String()
	if status != "completed" && status != "视频已缓存完成" && status != "" {
		logrus.Error("跳过未缓存完成的视频", status)
		return
	}
	return "video.m4s", "audio.m4s"
}

func OpenFolder(outputDir string) {
	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("explorer", outputDir).Start()
	case "darwin": // macOS
		_ = exec.Command("open", outputDir).Start()
	default: // Linux and other Unix-like systems
		_ = exec.Command("xdg-open", outputDir).Start()
	}
}
