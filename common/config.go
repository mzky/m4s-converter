package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"m4s-converter/internal"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/google/go-github/v65/github"
	"github.com/integrii/flaggy"
	"github.com/sirupsen/logrus"
)

// Config 存储应用程序的配置和状态
type Config struct {
	AssOFF     bool   // 关闭自动生成弹幕功能
	Overlay    bool   // 合成文件时是否覆盖同名视频
	Summarize  bool   // 将未合并的MP3和视频文件放入汇总目录
	CachePath  string // 自定义视频缓存路径
	GPACPath   string // 自定义GPAC的mp4box文件路径
	OutputDir  string // 输出目录
	Title      string // 视频标题
	Uname      string // UP主名称
	GroupTitle string // 分组标题
	ItemId     string // 项目ID
	GroupId    string // 分组ID
	Uid        string // 用户ID
	video      string // 当前处理的视频路径（内部使用）
	AssPath    string // 生成的ASS字幕路径（内部使用）
}

func (c *Config) flag() {
	var ver bool
	u, err := user.Current()
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false
	flaggy.SetName(color.CyanString("m4s-converter"))
	flaggy.SetDescription(color.CyanString("BiliBili音视频合成工具."))
	flaggy.Bool(&ver, "v", "version", "查看版本信息")
	flaggy.Bool(&c.AssOFF, "a", "assoff", "关闭自动生成弹幕功能，默认不关闭")
	flaggy.Bool(&c.Overlay, "o", "overlay", "合成文件时是否覆盖同名视频，默认不覆盖并重命名新文件")
	flaggy.Bool(&c.Summarize, "u", "summarize", "将未合并的MP3和视频文件放入汇总目录，默认不汇总")
	flaggy.String(&c.CachePath, "c", "cachepath", "自定义视频缓存路径，默认使用bilibili的默认缓存路径")
	flaggy.String(&c.GPACPath, "g", "gpacpath", "自定义GPAC的mp4box文件路径,值为select时弹出选择对话框")
	flaggy.ShowHelpOnUnexpectedEnable() // 解析到未预期参数时显示帮助
	flaggy.Parse()
	if ver {
		fmt.Println(color.CyanString("当前版本: %s", version))
		fmt.Println(color.CyanString("编译信息: %s", buildTime))
		fmt.Println(color.CyanString("源码版本: %s", sourceVer))
		os.Exit(0)
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
	if c.CachePath == "" {
		if err != nil {
			logrus.Warn("获取当前用户失败，使用默认缓存路径: ", err)
			c.CachePath = filepath.Join("~", "Videos", "bilibili")
		} else {
			c.CachePath = filepath.Join(u.HomeDir, "Videos", "bilibili")
		}
	}
	c.GetCachePath()
}
func (c *Config) InitConfig(ctx context.Context) {
	go c.PanicHandler()

	// 首先解析命令行参数
	c.flag()

	// 显示免责声明
	fmt.Println("=====================================================")
	fmt.Println("           使用本程序需遵守以下使用条款")
	fmt.Println("   仅转换本人通过哔哩哔哩官方客户端合法缓存的视频，")
	fmt.Println("  且转换结果严格用于个人备份，绝不传播、分享或商用。")
	fmt.Println("=====================================================")
	fmt.Println("  按任意键同意并继续使用，关闭窗口则拒绝并退出程序！")
	fmt.Println("=====================================================")

	// 等待用户输入任意键
	_, _ = fmt.Scanln()
	logrus.Info("用户同意使用，程序继续执行")

	diffVersion(ctx)
}

func diffVersion(ctx context.Context) {
	apiURL := "https://api.github.com/repos/mzky/m4s-converter/releases/latest"

	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 使用 context 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
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
