//go:build windows

package internal

import (
	_ "embed"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

//go:embed windows/MP4Box.exe
var mp4Box []byte

func GetMP4Box() string {
	mp4boxName := "MP4Box.exe"
	wd, _ := os.Getwd()
	mp4boxPath := filepath.Join(wd, mp4boxName) // 指定ffmpeg路径
	if !exist(mp4boxPath) {
		logrus.Info("第一次运行,自动释放MP4Box")
		if err := os.WriteFile(mp4boxName, mp4Box, os.ModePerm); err != nil {
			logrus.Error(err)
			logrus.Fatal("释放MP4Box失败,查看文件权限是否正常")
		}
	}
	return mp4boxPath
}
