package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"log"
	"m4s-converter/common"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defer panicHandler()
	var c common.Config
	c.InitConfig()

	// 查找m4s文件，并转换为mp4和mp3
	if err := filepath.WalkDir(c.CachePath, c.FindM4sFiles); err != nil {
		c.MessageBox(fmt.Sprintf("找不到 bilibili 的缓存 m4s 文件：%v", err))
		os.Exit(1)
	}

	dirs, err := common.GetCacheDir(c.CachePath)
	if err != nil {
		c.MessageBox(fmt.Sprintf("找不到 bilibili 的缓存目录：%v", err))
		os.Exit(1)
	}

	if dirs == nil {
		if common.Exist(filepath.Join(c.CachePath, ".videoInfo")) {
			dirs = append(dirs, c.CachePath)
		}
	}

	// 合成音视频文件
	var outputFiles []string
	for _, v := range dirs {
		video, audio, e := common.GetAudioAndVideo(v)
		if e != nil {
			log.Println("找不到已修复的音频和视频文件：", err)
			continue
		}
		info := filepath.Join(v, ".videoInfo")
		infoStr, _ := os.ReadFile(info)
		js, _ := simplejson.NewJson(infoStr)

		groupTitle, _ := js.Get("groupTitle").String()
		title, _ := js.Get("title").String()
		uname, _ := js.Get("uname").String()
		outputDir := filepath.Join(filepath.Dir(v), "output")
		if !common.Exist(outputDir) {
			_ = os.Mkdir(outputDir, os.ModePerm)
		}
		groupDir := filepath.Join(outputDir, groupTitle)
		if !common.Exist(groupDir) {
			if os.Mkdir(groupDir, os.ModePerm) != nil {
				c.MessageBox("无权限创建目录：" + groupDir)
				os.Exit(1)
			}
		}
		outputFile := filepath.Join(groupDir, title+"-"+uname+".mp4")
		if er := c.Composition(video, audio, outputFile); er != nil {
			log.Println(er)
			continue
		}
		outputFiles = append(outputFiles, outputFile)
	}

	if outputFiles != nil {
		log.Println("任务已全部完成:")
		fmt.Println(strings.Join(outputFiles, "\n"))
	}
	var input string
	fmt.Println("按回车键退出...")
	fmt.Scanln(&input)
}

func panicHandler() {
	if r := recover(); r != nil {
		fmt.Println("FFmpeg执行异常:", r)
		var input string
		fmt.Println("按回车键退出...")
		fmt.Scanln(&input)
	}
}
