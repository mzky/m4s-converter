package common

import (
	"fmt"
	"log"
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
		for {
			buf := make([]byte, 1024)
			n, e := stderr.Read(buf)
			if e != nil {
				return
			}
			cmdErr := string(buf[:n])
			//fmt.Print(cmdErr)
			if strings.Contains(cmdErr, "exists") {
				log.Println("视频文件已存在，跳过生成！")
				return
			}
		}
	}()

	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		return err
	}

	log.Println("已合成视频文件：\n", outputFile)
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
			dirs = append(dirs, path)
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
	// 读取源文件内容
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// 截取从第10个字符开始的数据,另存音视频文件
	err = os.WriteFile(dst, data[9:], os.ModePerm)
	if err != nil {
		return err
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
		log.Println("选择的 bilibili 缓存目录为: ", c.CachePath)
		return
	}
	c.MessageBox("未使用 bilibili 默认缓存路径 " + videosDir + " ，\n请选择 bilibili 当前设置的缓存路径！")
	c.SelectDirectory()
}

func (c *Config) GetFFmpegPath() {
	wd, _ := os.Getwd()
	c.FFmpegPath = filepath.Join(wd, "ffmpeg.exe") // 指定ffmpeg路径
	if !Exist(c.FFmpegPath) {
		c.SelectFile()
	}
}

func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
