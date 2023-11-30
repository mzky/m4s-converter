package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ffmpegPath = "c:\\ff\\ffmpeg.exe"                  // 指定ffmpeg路径
	cachePath  = "C:\\Users\\mzky\\Videos\\bilibili\\" // 指定要操作的目标文件路径
)

func main() {

	// 使用WalkDir遍历目录及其子目录
	if err := filepath.WalkDir(cachePath, findM4sFiles); err != nil {
		fmt.Println("Error walking the directory:", err)
		return
	}
	dirs, _ := getCacheDir(cachePath)
	for _, v := range dirs {
		files, _ := getAudioAndVideo(v)
		complex(files[0], files[1], filepath.Join(filepath.Dir(files[0]), "dddddd.mp4"))
	}
}

// C:\ff\ffmpeg.exe -i 1333045397-1-100050.m4s-video.mp4 -i 1333045397-1-30280.m4s-audio.mp3 -c:v copy -c:a aac -strict experimental out.mp4
func complex(videoFile, audioFile, outputFile string) error {
	// 构建FFmpeg命令行参数
	args := []string{
		"-i", videoFile,
		"-i", audioFile,
		"-c:v", "copy",
		"-c:a", "aac",
		"-strict", "experimental",
		"-y", // 是否替换已存在视频
		outputFile,
		"-hide_banner", // 隐藏版本信息和版权声明
		"-stats",       // 只显示统计信息
	}

	cmd := exec.Command(ffmpegPath, args...)

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

func findM4sFiles(src string, info os.DirEntry, err error) error {
	if err != nil {
		return err
	}
	// 查找.m4s文件
	if filepath.Ext(info.Name()) == ".m4s" {
		var dst string
		if strings.Contains(info.Name(), "30280") {
			dst = src + "-audio.mp3"
		} else {
			dst = src + "-video.mp4"
		}
		if err = m4sToAudioOrVideo(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func getCacheDir(cachePath string) ([]string, error) {
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

func getAudioAndVideo(cachePath string) ([]string, error) {
	var files []string
	err := filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.Contains(path, "audio.mp3") || strings.Contains(path, "video.mp4") {
				files = append(files, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func m4sToAudioOrVideo(src, dst string) error {
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
