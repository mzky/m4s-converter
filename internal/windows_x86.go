//go:build windows && !aarch

package internal

import (
	"embed"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

//go:embed ffmpeg.exe
var FFMpegFile embed.FS

var ffMpegName = "ffmpeg.exe"

func GetFFMpeg() string {
	wd, _ := os.Getwd()
	ffMpegPath := filepath.Join(wd, ffMpegName) // 指定ffmpeg路径
	if !exist(ffMpegPath) {
		logrus.Info("第一次运行,自动释放ffmpeg")
		if err := decFile(); err != nil {
			logrus.Error(err)
			return ffMpegName
		}
	}
	return ffMpegPath
}

// DecFile 解压ffmpeg.exe
func decFile() error {
	file, err := FFMpegFile.Open(ffMpegName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用文件
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return os.WriteFile(ffMpegName, data, os.ModePerm)
}

func exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
