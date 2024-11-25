package common

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/google/go-github/v65/github"
	"github.com/integrii/flaggy"
	"github.com/sirupsen/logrus"
	"io"
	"m4s-converter/internal"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

func (c *Config) flag() {
	var ver bool
	u, _ := user.Current()
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false
	flaggy.SetName(color.CyanString("m4s-converter"))
	flaggy.SetDescription(color.CyanString("BiliBili音视频合成工具."))
	flaggy.Bool(&ver, "v", "version", "查看版本信息")
	flaggy.Bool(&c.AssOFF, "a", "assoff", "关闭自动生成弹幕功能，默认不关闭")
	flaggy.Bool(&c.Skip, "s", "skip", "跳过合成已存在的同名视频（比overlay优先级高），默认不跳过")
	flaggy.Bool(&c.Overlay, "o", "overlay", "合成文件时是否覆盖同名视频，默认不覆盖并重命名新文件")
	flaggy.String(&c.CachePath, "c", "cachepath", "自定义视频缓存路径，默认使用bilibili的默认缓存路径")
	flaggy.String(&c.GPACPath, "g", "gpacpath", "自定义GPAC的mp4box文件路径,值为select时弹出选择对话框")
	flaggy.String(&c.FFMpegPath, "f", "ffmpegpath", "自定义FFMpeg文件路径,值为select时弹出选择对话框")
	flaggy.ShowHelpOnUnexpectedEnable() // 解析到未预期参数时显示帮助
	flaggy.Parse()
	if ver {
		fmt.Println(color.CyanString("当前版本: %s", version))
		fmt.Println(color.CyanString("编译信息: %s", buildTime))
		fmt.Println(color.CyanString("源码版本: %s", sourceVer))
		os.Exit(0)
	}
	if c.CachePath == "" {
		c.CachePath = filepath.Join(u.HomeDir, "Videos", "bilibili")
	}
	c.GetCachePath()
	if c.FFMpegPath != "" {
		if c.FFMpegPath == "select" {
			c.SelectFFMpegPath()
		}
		logrus.Warnln("使用FFMpeg进行音视频合成")
		return
	}
	if c.GPACPath != "" {
		if c.GPACPath == "select" {
			c.SelectGPACPath()
			logrus.Warnln("使用MP4Box进行音视频合成")
		}
		return
	}
	c.GPACPath = internal.GetMP4Box()
	logrus.Warnln("使用MP4Box进行音视频合成")
}
func (c *Config) InitConfig() {
	diffVersion()
	c.flag()
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
			logrus.Warnln("发现新版本:", latestVersion, fmt.Sprintf("(当前版本:%s)", version))
			logrus.Println("按住Ctrl并点击链接下载新版本:", releaseURL)
			fmt.Print("按[回车]跳过更新...")
			_, _ = fmt.Scanln()
		}
	}
}
