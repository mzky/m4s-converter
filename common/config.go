package common

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/google/go-github/v65/github"
	"github.com/gookit/goutil/cflag"
	"github.com/sirupsen/logrus"
	"io"
	"m4s-converter/internal"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func (c *Config) InitConfig() {
	u, _ := user.Current()
	f := cflag.New(func(cf *cflag.CFlags) {
		cf.Desc = "BiliBili synthesis tool."
		cf.Version = fmt.Sprintf("%s,%s,%s", version, sourceVer, buildTime)
	})
	f.BoolVar(&c.AssOFF, "assOFF", false, "是否关闭自动生成ass弹幕，默认不关闭;;a")
	f.StringVar(&c.FFMpegPath, "ffMpeg", "", "自定义FFMpeg文件路径;;f")
	f.StringVar(&c.CachePath, "cachePath", filepath.Join(u.HomeDir, "Videos", "bilibili"),
		"自定义缓存路径，默认使用BiliBili的默认路径;;c")
	overlay := f.Bool("overlay", false, "是否覆盖已存在的视频，默认不覆盖;;o")
	f.StringVar(&c.GPACPath, "gpacpath", "", "自定义GPAC的mp4box文件路径,替代FFMpeg合成文件\n参数为select时则弹出对话框选择文件;;g")
	help := f.Bool("help", false, "帮助信息;;h")
	_ = f.Parse(nil)
	if *help {
		f.ShowHelp()
		os.Exit(0)
	}

	diffVersion()
	if c.GPACPath != "" {
		if c.GPACPath == "select" {
			c.SelectGPACPath()
		}
	} else {
		if c.FFMpegPath == "" {
			c.FFMpegPath = internal.GetFFMpeg()
		}
	}
	c.GetCachePath()
	if *overlay {
		c.Overlay = "-y"
	} else {
		c.Overlay = "-n"
	}
}

func diffVersion() {
	apiURL := "https://api.github.com/repos/mzky/m4s-converter/releases/latest"
	resp, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var release *github.RepositoryRelease
	if json.Unmarshal(body, &release) != nil {
		return
	}

	// 解析版本号
	v, err := semver.NewVersion(version)
	if err != nil {
		return
	}

	latestVersion := release.GetTagName()
	lv, err := semver.NewVersion(latestVersion)
	if err != nil {
		return
	}

	releaseURL := fmt.Sprintf(
		"https://github.com/mzky/m4s-converter/releases/download/%s/%s", latestVersion, filepath.Base(os.Args[0]))
	// 版本号比较
	if !v.Equal(lv) {
		if v.LessThan(lv) {
			// MessageBox(fmt.Sprintf("发现新版本: %s\n访问 %s 下载新版本", latestVersion, releaseURL))
			logrus.Println("发现新版本:", latestVersion)
			logrus.Println("按住Ctrl并点击链接下载:", releaseURL)
			fmt.Print("按[回车]跳过更新...")
			_, _ = fmt.Scanln()
		}
	}
}
