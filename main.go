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
	common.InitLog()
	var c common.Config
	c.InitConfig()

	defer c.PanicHandler()
	defer c.File.Close()

	begin := time.Now().Unix()
	logrus.Println("查找缓存目录下可转换的文件...")
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
	c.OutputDir = filepath.Join(c.CachePath, "output")
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
			if !common.Exist(info) {
				info = filepath.Join(v, conver.PlayEntryJson)
				if !common.Exist(info) {
					continue
				}
			}
		}
		infoStr, e := os.ReadFile(info)
		if e != nil {
			logrus.Error("找不到包含视频信息的info相关文件: ", info)
			continue
		}
		js, e := simplejson.NewJson(infoStr)
		if e != nil {
			logrus.Error("videoInfo相关文件解析失败: ", info)
			continue
		}

		groupTitle := common.Filter(js.Get("groupTitle").String())
		if groupTitle == "" {
			groupTitle = common.Filter(js.Get("owner_name").String())
		}
		title := common.Filter(js.Get("page_data").Get("download_subtitle").String())
		if title == "" {
			title = common.Filter(js.Get("title").String())
		}
		uname := common.Filter(js.Get("uname").String())
		if uname == "" {
			uname = common.Filter(js.Get("title").String())
		}
		status := common.Filter(js.Get("status").String())
		if status == "" {
			status = common.Filter(js.Get("page_data").Get("download_title").String())
		}
		itemId, e := js.Get("itemId").Int()
		if itemId == 0 || e != nil {
			itemId, _ = js.Get("owner_id").Int()
		}
		if status != "completed" && status != "视频已缓存完成" {
			skipFilePaths = append(skipFilePaths, v)
			logrus.Warn("未缓存完成,跳过合成", v, title+"-"+uname)
			continue
		}
		if !common.Exist(c.OutputDir) {
			_ = os.Mkdir(c.OutputDir, os.ModePerm)
		}
		groupDir := filepath.Join(c.OutputDir, groupTitle+"-"+uname)
		if !common.Exist(groupDir) {
			if err = os.Mkdir(groupDir, os.ModePerm); err != nil {
				common.MessageBox("无法创建目录：" + groupDir)
				wait()
			}
		}
		outputFile := filepath.Join(groupDir, title+conver.Mp4Suffix)
		if common.Exist(outputFile) && c.Overlay == "-n" {
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
		logrus.Printf("合成的文件:\n%s\n输出目录: %s",
			strings.ReplaceAll(strings.Join(outputFiles, "\n"), c.OutputDir, ""),
			c.OutputDir)
		// 打开合成文件目录
		go common.OpenFolder(c.OutputDir)
	} else {
		logrus.Warn("未合成任何文件！")
	}
	logrus.Print("已完成本次任务，耗时:", end-begin, "秒")
	logrus.Print("==========================================")

	wait()
}
func wait() {
	fmt.Print("按[回车]键退出...")
	_, _ = fmt.Scanln()
	os.Exit(0)
}
