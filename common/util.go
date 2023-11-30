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

var (
	wd, _      = os.Getwd()
	FFmpegPath = filepath.Join(wd, "ffmpeg.exe") // 指定ffmpeg路径
	CachePath  = GetCachePath()                  // 指定要操作的目标文件路径
	Overlay    = "-n"
)

func Composition(videoFile, audioFile, outputFile string) error {
	// 构建FFmpeg命令行参数
	args := []string{
		"-i", videoFile,
		"-i", audioFile,
		"-c:v", "copy",
		"-c:a", "aac",
		"-strict", "experimental",
		Overlay, // 是否覆盖已存在视频
		outputFile,
		"-hide_banner", // 隐藏版本信息和版权声明
		"-stats",       // 只显示统计信息
	}

	cmd := exec.Command(FFmpegPath, args...)

	// 设置输出和错误流 pipe
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// 启动命令
	err := cmd.Start()
	if err != nil {
		fmt.Printf("启动命令失败: %s\n", err)
		os.Exit(1)
	}

	// 读取并打印输出流
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := stdout.Read(buf)
			if err != nil {
				return
			}
			fmt.Print(string(buf[:n]))
		}
	}()

	// 读取并打印错误流
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := stderr.Read(buf)
			if err != nil {
				return
			}
			cmdErr := string(buf[:n])
			fmt.Print(cmdErr)
			if strings.Contains(cmdErr, "exists") {
				fmt.Println("视频文件已存在，跳过合成！")
				return
			}
		}
	}()

	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		return err
	}

	log.Println("已合成视频文件：", outputFile)
	return nil
}

func FindM4sFiles(src string, info os.DirEntry, err error) error {
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
	err = os.WriteFile(dst, data[9:], 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetCachePath() string {
	user, err := user.Current()
	if err != nil {
		log.Println("无法获取当前用户：", err)
		return ""
	}

	videosDir := filepath.Join(user.HomeDir, "Videos", "bilibili")
	ext, err := exists(videosDir)
	if err != nil && ext {
		log.Println("检查目录是否存在:", err)
		return ""
	}

	return videosDir

}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
