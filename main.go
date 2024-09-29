package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/sirupsen/logrus"
	"m4s-converter/common"
	"m4s-converter/conver"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	var c common.Config
	c.InitConfig()

	defer c.PanicHandler()
	defer c.File.Close()

	begin := time.Now().Unix()
	fmt.Println("查找缓存目录下可转换的文件...")
	// 查找m4s文件，并转换为mp4和mp3
	if err := filepath.WalkDir(c.CachePath, c.FindM4sFiles); err != nil {
		common.MessageBox(fmt.Sprintf("找不到 BiliBili 目录下的 m4s 文件：%v", err))
		wait()
	}

	dirs, err := common.GetCacheDir(c.CachePath) // 缓存根目录模式
	if err != nil {
		common.MessageBox(fmt.Sprintf("找不到 BiliBili 的缓存目录：%v", err))
		wait()
	}

	if dirs == nil {
		// 判断非缓存根目录时，验证是否为子目录
		if common.Exist(filepath.Join(c.CachePath, conver.VideoInfoSuffix)) ||
			common.Exist(filepath.Join(c.CachePath, conver.VideoInfoJson)) {
			dirs = append(dirs, c.CachePath)
		}
	}

	// 合成音视频文件
	var outputDir string
	var outputFiles []string
	var skipFilePaths []string
	for _, v := range dirs {
		video, audio, e := c.GetAudioAndVideo(v)
		if e != nil {
			logrus.Error("找不到已修复的音频和视频文件:", err)
			continue
		}
		info := filepath.Join(v, conver.VideoInfoJson)
		if !common.Exist(info) {
			info = filepath.Join(v, conver.VideoInfoSuffix)
		}
		infoStr, e := os.ReadFile(info)
		if e != nil {
			logrus.Error("找不到videoInfo相关文件: ", info)
			continue
		}
		js, e := simplejson.NewJson(infoStr)
		if e != nil {
			logrus.Error("videoInfo相关文件解析失败: ", info)
			continue
		}

		groupTitle := common.Filter(js.Get("groupTitle").String())
		title := common.Filter(js.Get("title").String())
		uname := common.Filter(js.Get("uname").String())
		status := common.Filter(js.Get("status").String())
		itemId, _ := js.Get("itemId").Int()

		if status != "completed" {
			skipFilePaths = append(skipFilePaths, v)
			logrus.Warn("未缓存完成,跳过合成", v, title+"-"+uname)
			continue
		}
		outputDir = filepath.Join(filepath.Dir(v), "output")
		if !common.Exist(outputDir) {
			_ = os.Mkdir(outputDir, os.ModePerm)
		}
		groupDir := filepath.Join(outputDir, groupTitle+"-"+uname)
		if !common.Exist(groupDir) {
			if err = os.Mkdir(groupDir, os.ModePerm); err != nil {
				common.MessageBox("无法创建目录：" + groupDir)
				wait()
			}
		}
		outputFile := filepath.Join(groupDir, title+conver.Mp4Suffix)
		if common.Exist(outputFile) {
			outputFile = filepath.Join(groupDir, title+strconv.Itoa(itemId)+conver.Mp4Suffix)
		}
		if er := c.Composition(video, audio, outputFile); er != nil {
			logrus.Error("合成失败:", er)
			continue
		}
		outputFiles = append(outputFiles, outputFile)
	}

	end := time.Now().Unix()
	logrus.Print("==========================================")
	if skipFilePaths != nil {
		logrus.Print("跳过的目录:\n" + strings.Join(skipFilePaths, "\n"))
	}
	if outputFiles != nil {
		logrus.Print("合成的文件:\n" + strings.Join(outputFiles, "\n"))
		// 打开合成文件目录
		go common.OpenFolder(outputDir)
	} else {
		logrus.Warn("未合成任何文件！")
	}
	logrus.Print("已完成本次任务，耗时:", end-begin, "秒")
	logrus.Print("==========================================")

	wait()
}

func wait() {
	fmt.Print("按回车键退出...")
	_, _ = fmt.Scanln()
	os.Exit(0)
}
