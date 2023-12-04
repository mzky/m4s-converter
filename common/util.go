package common

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	FFmpegPath string
	CachePath  string
	Overlay    string
	File       *os.File
}

func (c *Config) InitConfig() {
	c.Flags()
	c.GetFFmpegPath()
	c.GetCachePath()
	c.Overlay = "-n"
}

func (c *Config) Composition(videoFile, audioFile, outputFile string) error {
	// 构建FFmpeg命令行参数
	args := []string{
		"-i", videoFile,
		"-i", audioFile,
		"-c:v", "copy", // video不指定编解码，使用bilibili原有编码
		"-c:a", "copy", // audio不指定编解码，使用bilibili原有编码
		"-strict", "experimental", // 宽松编码控制器
		c.Overlay, // 是否覆盖已存在视频
		outputFile,
		"-hide_banner", // 隐藏版本信息和版权声明
		"-stats",       // 只显示统计信息
	}

	cmd := exec.Command(c.FFmpegPath, args...)

	// 设置输出和错误流 pipe
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// 启动命令
	err := cmd.Start()
	if err != nil {
		c.MessageBox(fmt.Sprintf("执行FFmpeg命令失败: %s", err))
		os.Exit(1)
	}

	// 读取并打印输出流
	go func() {
		for {
			buf := make([]byte, 1024)
			n, e := stdout.Read(buf)
			if e != nil {
				return
			}
			fmt.Print(string(buf[:n]))
		}
	}()

	// 读取并打印错误流
	go func() {
		logrus.Print("视频生成中")
		for {
			buf := make([]byte, 1024)
			n, e := stderr.Read(buf)
			if e != nil {
				return
			}
			cmdErr := string(buf[:n])
			fmt.Print(".")
			if strings.Contains(cmdErr, "exists") {
				fmt.Println()
				logrus.Warn("跳过已经生成的音视频文件！")
			}
		}
	}()

	// 等待命令执行完成
	cmd.Wait()
	logrus.Info("已合成视频文件:", filepath.Base(outputFile))
	return nil
}

func (c *Config) FindM4sFiles(src string, info os.DirEntry, err error) error {
	if err != nil {
		return err
	}
	// 查找.m4s文件
	if filepath.Ext(info.Name()) == ".m4s" {
		var dst string
		if strings.Contains(info.Name(), "30280") { // 30280是音频文件
			dst = src + "-audio.mp3"
		} else {
			dst = src + "-video.mp4"
		}
		if err = M4sToAudioOrVideo(src, dst); err != nil {
			c.MessageBox(fmt.Sprintf("%v 转换异常：%v", src, err))
			return err
		}
		logrus.Info("已将m4s转换为音视频文件:", dst)
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

func GetAudioAndVideo(cachePath string) (string, string, error) {
	var video string
	var audio string
	err := filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.Contains(path, "video.mp4") {
				video = path
			}
			if strings.Contains(path, "audio.mp3") {
				audio = path
			}
		}
		return nil
	})

	if err != nil {
		return "", "", err
	}

	return video, audio, nil
}

func M4sToAudioOrVideo(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer srcFile.Close()

	// 读取前9个字符
	data := make([]byte, 9)
	io.ReadAtLeast(srcFile, data, 9)
	if string(data) != "000000000" {
		return fmt.Errorf("音视频文件不是9个0的头，跳过转换")
	}

	// 移动到第9个字节
	_, err = srcFile.Seek(9, 0) // 从文件开头偏移
	if err != nil {
		return fmt.Errorf("文件字节偏移失败: %v", err)
	}

	// 创建新文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建新文件失败: %v", err)
	}
	defer dstFile.Close()

	// 将截取后的内容写入新文件
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

func (c *Config) GetCachePath() {
	u, err := user.Current()
	if err != nil {
		c.MessageBox(fmt.Sprintf("无法获取当前用户：%v", err))
		return
	}

	videosDir := filepath.Join(u.HomeDir, "Videos", "bilibili")
	if Exist(videosDir) {
		c.CachePath = videosDir
		return
	}
	if Exist(filepath.Join(c.CachePath, ".videoInfo")) || Exist(filepath.Join(c.CachePath, "load_log")) {
		logrus.Info("选择的 bilibili 缓存目录为: ", c.CachePath)
		return
	}
	c.MessageBox("未使用 bilibili 默认缓存路径 " + videosDir + " ，请选择 bilibili 当前设置的缓存路径！")
	c.SelectDirectory()
}

func (c *Config) GetFFmpegPath() {
	wd, _ := os.Getwd()
	c.FFmpegPath = filepath.Join(wd, "ffmpeg.exe") // 指定ffmpeg路径
	if !Exist(c.FFmpegPath) {
		logrus.Info("找不到ffmpeg.exe,自动下载ffmpeg...")
		if err := DownloadFile(); err != nil {
			logrus.Error(err)
			return
		}
		if !Exist(c.FFmpegPath) {
			logrus.Warn("无法自动下载ffmpeg,打开本地ffmpeg.exe文件")
			c.SelectFile()
		}
	}
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func DownloadFile() error {
	url := "https://mirror.ghproxy.com/https://github.com/mzky/m4s-converter/releases/download/ffmpeg/ffmpeg.exe"
	filename := "ffmpeg.exe"

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("无法访问: %v", err)
	}
	defer response.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("无权创建: %v", err)
	}
	defer file.Close()

	bar := pb.Full.Start64(response.ContentLength)

	reader := bar.NewProxyReader(response.Body)
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("写入失败: %v", err)
	}

	bar.Finish()
	return nil
}

func Filter(name string, err error) string {
	name = strings.ReplaceAll(name, "<", "《")
	name = strings.ReplaceAll(name, ">", "》")
	name = strings.ReplaceAll(name, `\`, "#")
	name = strings.ReplaceAll(name, `"`, "'")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "|", "_")
	name = strings.ReplaceAll(name, "?", "_")
	name = strings.ReplaceAll(name, "*", "_")

	return name
}
