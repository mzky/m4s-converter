package common

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

//go:embed ffmpeg.exe
var ffmpegFile embed.FS

var (
	FFmpegName    = "ffmpeg.exe"
	FileHashValue = "e3de8aad89e68d2f161050fb97a6568a2d8ff3ca0eae695448097e4d174a02d1"
)

type Config struct {
	FFmpegPath string
	CachePath  string
	Overlay    string
	File       *os.File
}

func (c *Config) InitConfig() {
	InitLog()
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
	if err := cmd.Start(); err != nil {
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
		fmt.Print("准备合成mp4 ...")
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
				logrus.Warn("跳过已经存在的音视频文件:", filepath.Base(outputFile))
			}
		}
	}()

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
	c.FFmpegPath = filepath.Join(wd, FFmpegName) // 指定ffmpeg路径
	if !Exist(c.FFmpegPath) {
		logrus.Info("第一次运行,自动释放ffmpeg.exe")
		if err := DecFile(); err != nil {
			logrus.Error(err)
		}
	}
	if !c.FileHashCompare() {
		logrus.Info("文件不完整,重新释放ffmpeg.exe")
		if err := DecFile(); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func DecFile() error {
	file, err := ffmpegFile.Open(FFmpegName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用文件
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return os.WriteFile(FFmpegName, data, os.ModePerm)
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
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

func (c *Config) PanicHandler() {
	if e := recover(); e != nil {
		c.File.Close()
		fmt.Print("按回车键退出...")
		fmt.Scanln()
	}
}

func (c *Config) FileHashCompare() bool {
	file, err := os.ReadFile(c.FFmpegPath)
	if err != nil {
		logrus.Error("打开文件失败:", err)
		return false
	}

	// 计算文件的SHA-256哈希值
	hash := sha256.Sum256(file)
	sha256Str := fmt.Sprintf("%x", hash)

	return FileHashValue == sha256Str
}
