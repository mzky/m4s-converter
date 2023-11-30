package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"log"
	"m4s-converter/common"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	ffmpegPath := flag.String("f", common.FFmpegPath, "指定ffmpeg路径，或将本程序访问ffmpeg.exe同目录")
	cachePath := flag.String("c", common.CachePath, "指定bilibili缓存目录")
	overlay := flag.Bool("y", false, "是否覆盖，默认不覆盖")

	flag.Parse()

	if *overlay {
		common.Overlay = "-y"
	}
	if *ffmpegPath != common.FFmpegPath {
		common.FFmpegPath = *ffmpegPath
	}
	if *cachePath != common.CachePath {
		common.CachePath = *cachePath
	}

	// 使用WalkDir遍历目录及其子目录
	if err := filepath.WalkDir(common.CachePath, common.FindM4sFiles); err != nil {
		log.Println("找不到bilibili的缓存m4s文件:", err)
		return
	}
	dirs, err := common.GetCacheDir(common.CachePath)
	if err != nil {
		log.Println("找不到bilibili的缓存目录:", err)
		return
	}
	var outputFiles []string
	for _, v := range dirs {
		video, audio, err := common.GetAudioAndVideo(v)
		if err != nil {
			log.Println("找不到音频和视频文件:", err)
			continue
		}
		info := filepath.Join(v, ".videoInfo")
		infoStr, _ := os.ReadFile(info)
		js, _ := simplejson.NewJson(infoStr)

		//groupTitle, _ := js.Get("groupTitle").String()
		title, _ := js.Get("title").String()
		uname, _ := js.Get("uname").String()
		outputFile := filepath.Join(v, title+"-"+uname+".mp4")
		if err := common.Composition(video, audio, outputFile); err != nil {
			log.Println(err)
			continue
		}
		outputFiles = append(outputFiles, outputFile)
	}
	log.Println("任务已全部完成:")
	fmt.Println(strings.Join(outputFiles, "\n"))

}
