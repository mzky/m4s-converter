package common

import (
	"encoding/json"
	"fmt"
	"github.com/ncruces/zenity"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"m4s-converter/conver"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	FFMpegPath string
	CachePath  string
	Overlay    string
	File       os.File
	AssPath    string
	AssOFF     bool
	OutputDir  string
	GPACPath   string
}

func (c *Config) Composition(videoFile, audioFile, outputFile string) error {
	var cmd *exec.Cmd
	if c.GPACPath != "" {
		cmd = exec.Command(c.GPACPath, "-add", videoFile, "-add", audioFile, outputFile)
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
			c.Overlay,                 // 是否覆盖已存在视频
			"-movflags", "+faststart", // 启用faststart可以让视频在网络传输时更快地开始播放
			outputFile,
			"-hide_banner", // 隐藏版本信息和版权声明
			"-stats",       // 只显示统计信息
		)
		// fmt.Println("执行FFmpeg命令:", c.FFMpegPath, strings.Join(args, " "))
		// logrus.Info(c.FFMpegPath, args)
		cmd = exec.Command(c.FFMpegPath, args...)
	}

	// 设置输出和错误流 pipe
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// 启动命令
	if err := cmd.Start(); err != nil {
		logrus.Error(err)
		MessageBox("执行FFmpeg命令失败,查看程序是否有足够权限和是否安装了FFmpeg")
		os.Exit(1)
	}

	// 读取并打印输出流
	go printOutput(stdout)

	// 读取并打印错误流
	go printError(stderr, outputFile)

	assFile := strings.ReplaceAll(outputFile, filepath.Ext(outputFile), conver.AssSuffix)
	if !Exist(assFile) {
		logrus.Error("下载弹幕文件失败！")
	} else {
		_ = c.copyFile(c.AssPath, assFile)
	}

	// 等待命令执行完成
	if err := cmd.Wait(); err == nil {
		logrus.Info("已合成视频文件:", filepath.Base(outputFile))
	}
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
		logrus.Info("已将m4s转换为音视频文件:\n", dst)
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
	// return "https://api.bilibili.com/x/v1/dm/list.so?oid=" + cid
	return "https://comment.bilibili.com/" + cid + conver.XmlSuffix
}

// GetAudioAndVideo 从给定的缓存路径中查找音频和视频文件，并尝试下载并转换xml弹幕为ass格式
// 参数:
// - cachePath: 缓存路径，用于搜索音频、视频文件以及存储下载的弹幕文件
// 返回值:
// - video: 查找到的视频文件路径
// - audio: 查找到的音频文件路径
// - error: 在搜索、下载或转换过程中遇到的任何错误
func (c *Config) GetAudioAndVideo(cachePath string) (string, string, error) {
	var video string
	var audio string

	// 遍历给定路径下的所有文件和目录
	err := filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 如果遇到错误，立即返回
		}
		if !info.IsDir() {
			// 如果是文件，检查是否为视频或音频文件
			if strings.Contains(path, conver.VideoSuffix) {
				video = path // 找到视频文件
			}
			if strings.Contains(path, conver.AudioSuffix) {
				audio = path // 找到音频文件
			}
		} else {
			// 如果是目录，尝试下载并转换xml弹幕为ass格式
			if !c.AssOFF {
				danmakuXml := filepath.Join(path, conver.DanmakuXml)
				if Exist(danmakuXml) {
					c.AssPath = conver.Xml2ass(danmakuXml) // 转换xml弹幕文件为ass格式
					return nil
				}
				xmlPath := filepath.Join(path, info.Name()+conver.XmlSuffix)
				if e := downloadFile(joinUrl(info.Name()), xmlPath); e != nil {
					return nil
				}
				c.AssPath = conver.Xml2ass(xmlPath) // 转换xml弹幕文件为ass格式
			}
		}
		return nil
	})

	if err != nil {
		return "", "", err // 如果遍历过程中发生错误，返回错误信息
	}

	return video, audio, nil // 返回找到的视频和音频文件路径
}

func (c *Config) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		logrus.Error("文件打开失败:", err)
		return err
	}
	defer srcFile.Close()

	if _, err := srcFile.Seek(0, 0); err != nil {
		logrus.Errorf("文件指针重置失败: %v", err)
		return err
	}
	// 读取前9个字符
	data := make([]byte, 9)
	if _, err := io.ReadAtLeast(srcFile, data, 9); err != nil {
		logrus.Errorf("读取文件头失败: %v", err)
		return err
	}
	if string(data) == "000000000" {
		// 移动到第9个字节
		_, err = srcFile.Seek(9, 0) // 从文件开头偏移
		if err != nil {
			logrus.Errorf("文件字节偏移失败: %v", err)
			return err
		}
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	_ = srcFile.Sync()
	_ = dstFile.Sync()
	return nil
}

func (c *Config) M4sToAV(src, dst string) error {
	return c.copyFile(src, dst)
}

// GetCachePath 获取用户视频缓存路径
func (c *Config) GetCachePath() {
	if c.findM4sFiles() != nil {
		MessageBox("BiliBili缓存路径 " + c.CachePath + " 未找到缓存文件,\n请重新选择 BiliBili 缓存文件路径！")
		c.SelectDirectory()
		return
	}
	logrus.Info("选择的 BiliBili 缓存目录为: ", c.CachePath)
	return
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// Filter 过滤文件名
func Filter(name string, err error) string {
	if err != nil || name == "" {
		return ""
	}
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
		_ = c.File.Close()
		fmt.Print("按回车键退出...")
		_, _ = fmt.Scanln()
	}
}

func MessageBox(text string) {
	logrus.Warn(text)
	_ = zenity.Warning(text, zenity.Title("提示"), zenity.Width(400))
}

// checkFilesExist 检查多个文件是否存在
func checkFilesExist(paths ...string) bool {
	for _, path := range paths {
		if Exist(path) {
			return true
		}
	}
	return false
}

// findM4sFiles 检查目录及其子目录下是否存在m4s文件
func (c *Config) findM4sFiles() error {
	var m4sFiles []string
	err := filepath.Walk(c.CachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Error("遍历目录异常: %v, 文件路径: ", err, path)
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

	if checkFilesExist(c.GPACPath) {
		logrus.Info("选择 GPAC 的 mp4box 文件为:", c.CachePath)
		return
	}
	MessageBox("选择 GPAC 的 mp4box 文件不存在，请重新选择！")
	c.SelectDirectory()
}

func printOutput(stdout io.ReadCloser) {
	buf := make([]byte, 1024)
	for {
		n, e := stdout.Read(buf)
		if e != nil {
			// logrus.Error("读取标准输出错误:", e)
			return
		}
		if n > 0 {
			fmt.Print(string(buf[:n]))
		}
	}
}

func printError(stderr io.ReadCloser, outputFile string) {
	logrus.Println("# 准备合成:", filepath.Base(outputFile))
	buf := make([]byte, 1024)
	for {
		n, e := stderr.Read(buf)
		if e != nil {
			// logrus.Error("读取标准错误输出错误:", e)
			return
		}
		if n > 0 {
			cmdErr := string(buf[:n])
			if strings.Contains(cmdErr, "exists") {
				logrus.Warn("跳过已经存在的音视频文件:", filepath.Base(outputFile))
			}
		}
	}
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
		puStr := string(puByte)
		var p gjson.Result
		if p = gjson.Get(puStr, "data"); !p.Exists() {
			p = gjson.Get(puStr, "result")
		}
		if p.Exists() {
			return p.Get("dash.video|@reverse|0.id").String(), p.Get("dash.audio|@reverse|0.id").String()
		}
		return "", ""
	}
	logrus.Warn("找不到.playurl文件:\n", pu)
	pu = filepath.Join(filepath.Dir(filepath.Dir(patch)), conver.PlayEntryJson)
	puDate, e := os.ReadFile(pu)
	if e != nil {
		logrus.Error("找不到entry.json文件: ", pu)
		return
	}
	var p conver.Entry
	if err := json.Unmarshal(puDate, &p); err != nil {
		logrus.Error("解析entry.json文件失败: ", err)
		return
	}
	if p.PageData.DownloadTitle != "视频已缓存完成" {
		logrus.Error("跳过未缓存完成的视频", p.PageData.DownloadSubtitle)
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
