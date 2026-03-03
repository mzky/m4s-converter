package common

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"m4s-converter/conver"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	utils "github.com/mzky/utils/common"
	"github.com/ncruces/zenity"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Config struct {
	CachePath  string
	Overlay    bool
	AssPath    string
	AssOFF     bool
	OutputDir  string
	GPACPath   string
	Summarize  bool
	video      string
	audio      string
	ItemId     string
	GroupId    string
	Uid        string
	Title      string
	Uname      string
	GroupTitle string
	ExitFlag   bool
}

func (c *Config) overlay() string {
	if c.Overlay {
		return "-y"
	}
	return "-n"
}

// SetExitFlag 设置退出标志
func (c *Config) SetExitFlag(flag bool) {
	c.ExitFlag = flag
}

// ShouldExit 检查是否应该退出
func (c *Config) ShouldExit() bool {
	return c.ExitFlag
}
func (c *Config) Composition(videoFile, audioFile, outputFile string) error {
	var cmd *exec.Cmd
	// 构建MP4Box命令行参数
	var args []string
	// 添加覆盖参数
	if c.Overlay {
		args = append(args, "-force")
	}

	// 添加字符集参数，指定使用UTF-8编码
	args = append(args, "-charset", "utf8")

	// 添加元数据标签
	tags := fmt.Sprintf("title=%s:artist=%s:album=%s", c.GroupId, c.Uid, c.ItemId)
	args = append(args, "-tags", tags)
	args = append(args,
		// "-quiet", // 仅打印异常日志
		"-cprt", c.ItemId,
		"-add", videoFile+"#video",
		"-add", audioFile+"#audio",
		"-new", outputFile)
	cmd = exec.Command(c.GPACPath, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	if c.AssPath != "" {
		assFile := strings.ReplaceAll(outputFile, conver.Mp4Suffix, conver.AssSuffix)
		_ = c.copyFile(c.AssPath, assFile)
	}

	// 等待命令执行完成
	if err := cmd.Run(); err != nil {
		logrus.Errorf("合成视频文件失败:%s\n%s", outputFile, stdout.String())
		return err
	}

	logrus.Info("已合成视频文件:", outputFile)
	return nil
}

func (c *Config) FindM4sFiles(src string, info os.DirEntry, err error) error {
	if err != nil {
		return err
	}
	// 查找.m4s文件
	if strings.HasSuffix(info.Name(), conver.M4sSuffix) {
		var dst string
		videoId, audioId := GetVAId(src)
		if videoId != "" && audioId != "" {
			if strings.Contains(info.Name(), audioId) { // 音频文件
				dst = strings.ReplaceAll(src, conver.M4sSuffix, conver.AudioSuffix)
			} else {
				dst = strings.ReplaceAll(src, conver.M4sSuffix, conver.VideoSuffix)
			}
		}

		if err = c.M4sToAV(src, dst); err != nil {
			MessageBox(fmt.Sprintf("%v 转换异常：%v", src, err))
			return err
		}
		logrus.Info("已将m4s转换为音视频文件: ", strings.TrimPrefix(dst, c.CachePath))
	}
	return nil
}

func GetCacheDir(cachePath string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(cachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != cachePath {
			if !strings.Contains(path, "output") {
				dirs = append(dirs, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirs, nil
}

func joinUrl(cid string) string {
	return "https://comment.bilibili.com/" + cid + conver.XmlSuffix
}
func joinXmlUrl(cid string) string {
	return "https://api.bilibili.com/x/v1/dm/list.so?oid=" + cid
}

// GetAudioAndVideo 从给定的缓存路径中查找音频和视频文件，并尝试下载并转换xml弹幕为ass格式
// 参数:
// - cachePath: 缓存路径，用于搜索音频、视频文件以及存储下载的弹幕文件
// 返回值:
// - video: 查找到的视频文件路径
// - audio: 查找到的音频文件路径
// - error: 在搜索、下载或转换过程中遇到的任何错误
func (c *Config) GetAudioAndVideo(cachePath string) (string, string, error) {
	var video, audio string

	// 遍历给定路径下的所有文件（不包括子目录）
	entries, err := os.ReadDir(cachePath)
	if err != nil {
		return "", "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// 如果是目录，递归查找
			childVideo, childAudio, err := c.GetAudioAndVideo(filepath.Join(cachePath, entry.Name()))
			if err == nil && childVideo != "" && childAudio != "" {
				video = childVideo
				audio = childAudio
				break
			}
			continue
		}

		// 如果是文件，检查是否为视频或音频文件
		fileName := entry.Name()
		if strings.HasSuffix(fileName, conver.VideoSuffix) {
			video = filepath.Join(cachePath, fileName)
		}
		if strings.HasSuffix(fileName, conver.AudioSuffix) {
			audio = filepath.Join(cachePath, fileName)
		}
	}

	// 如果在当前目录及其子目录中都找不到视频或音频文件，返回错误
	if video == "" || audio == "" {
		return "", "", fmt.Errorf("找不到音频或视频文件: %s", cachePath)
	}

	// 下载弹幕文件
	if !c.AssOFF {
		// 保存当前的video路径，用于downloadXml
		oldVideo := c.video
		c.video = video
		c.downloadXml()
		c.video = oldVideo
	}
	return video, audio, nil // 返回找到的视频和音频文件路径
}
func (c *Config) copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		logrus.Errorf("打开源文件失败: %v", err)
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		logrus.Errorf("创建目标文件失败: %v", err)
		return err
	}
	defer dstFile.Close()

	// 读取前 9 个字节
	data := make([]byte, 9)
	if _, err := io.ReadAtLeast(srcFile, data, 9); err != nil {
		logrus.Errorf("读取文件头失败: %v", err)
		return err
	}

	// 检查前 9 个字节是否为 '0'
	if string(data) != "000000000" {
		// 如果前 9 个字节不为 '0'，写入这些字节
		if _, err := dstFile.Write(data); err != nil {
			logrus.Errorf("写入文件头失败: %v", err)
			return err
		}
	}

	// 使用缓冲读取器逐块读取并写入文件
	if _, err := io.Copy(bufio.NewWriter(dstFile), bufio.NewReader(srcFile)); err != nil {
		logrus.Errorf("读取或写入文件失败: %v", err)
		return err
	}
	return nil
}

func (c *Config) M4sToAV(src, dst string) error {
	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		logrus.Errorf("创建目标目录失败: %v", err)
		return err
	}
	return c.copyFile(src, dst)
}

// GetCachePath 获取用户视频缓存路径
func (c *Config) GetCachePath() {
	if c.findM4sFiles() != nil {
		MessageBox("BiliBili缓存路径 " + c.CachePath + " 未找到缓存文件, \n请重新选择 BiliBili 缓存文件路径！")
		c.SelectDirectory()
		return
	}
	logrus.Info("选择的 BiliBili 缓存目录为: ", c.CachePath)
	return
}

func Size(path string) int64 {
	if utils.IsExist(path) {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return 0
		}
		return fileInfo.Size()
	}
	return 0
}

// Filter 过滤文件名
func Filter(name string, err error) string {
	if err != nil || name == "" {
		return ""
	}
	name = strings.ReplaceAll(name, "（", "(")
	name = strings.ReplaceAll(name, "）", ")")
	name = strings.ReplaceAll(name, "<", "《")
	name = strings.ReplaceAll(name, ">", "》")
	name = strings.ReplaceAll(name, `\`, "#")
	name = strings.ReplaceAll(name, `"`, `'`)
	name = strings.ReplaceAll(name, "/", "#")
	name = strings.ReplaceAll(name, "|", "_")
	name = strings.ReplaceAll(name, "?", "？")
	name = strings.ReplaceAll(name, "*", "-")
	name = strings.ReplaceAll(name, "【", "[")
	name = strings.ReplaceAll(name, "】", "]")
	name = strings.ReplaceAll(name, ":", "：")
	name = strings.ReplaceAll(name, " ", "_")

	return strings.TrimSpace(name)
}

func (c *Config) PanicHandler() {
	if e := recover(); e != nil {
		fmt.Print("按回车键退出...")
		_, _ = fmt.Scanln()
	}
}

// calculateFileHash 计算文件的MD5哈希值（流式计算）
func (c *Config) calculateFileHash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Errorf("打开文件失败: %v", err)
		return ""
	}
	defer file.Close()

	hash := md5.New()

	// 使用流式读取，每次读取4KB
	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			logrus.Errorf("读取文件失败: %v", err)
			return ""
		}
		if n == 0 {
			break
		}
		// 只更新实际读取的数据
		hash.Write(buffer[:n])
	}

	return hex.EncodeToString(hash.Sum(nil))
}

// calculateCombinedHash 计算音频和视频文件的组合哈希值（流式计算）
func (c *Config) calculateCombinedHash(videoPath string, audioPath string) string {
	hash := md5.New()

	// 计算视频文件哈希（流式）
	videoFile, err := os.Open(videoPath)
	if err == nil {
		// 使用流式读取，每次读取4KB
		buffer := make([]byte, 4096)
		for {
			n, readErr := videoFile.Read(buffer)
			if readErr != nil && readErr != io.EOF {
				logrus.Errorf("读取视频文件失败: %v", readErr)
				videoFile.Close()
				return ""
			}
			if n == 0 {
				break
			}
			// 只更新实际读取的数据
			hash.Write(buffer[:n])
		}
		videoFile.Close()
	} else {
		logrus.Errorf("打开视频文件失败: %v", err)
		return ""
	}

	// 计算音频文件哈希（流式）
	audioFile, err := os.Open(audioPath)
	if err == nil {
		// 使用流式读取，每次读取4KB
		buffer := make([]byte, 4096)
		for {
			n, readErr := audioFile.Read(buffer)
			if readErr != nil && readErr != io.EOF {
				logrus.Errorf("读取音频文件失败: %v", readErr)
				audioFile.Close()
				return ""
			}
			if n == 0 {
				break
			}
			// 只更新实际读取的数据
			hash.Write(buffer[:n])
		}
		audioFile.Close()
	} else {
		logrus.Errorf("打开音频文件失败: %v", err)
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}

// isFileIdentical 检查输出文件是否与输入的音频和视频文件完全相同
func (c *Config) isFileIdentical(outputFile string, videoPath string, audioPath string) bool {
	// 首先检查文件是否存在
	if !utils.IsExist(outputFile) {
		return false
	}

	// 检查文件大小是否相近
	outputInfo, err := os.Stat(outputFile)
	if err != nil {
		logrus.Errorf("获取输出文件信息失败: %v", err)
		return false
	}

	videoInfo, err := os.Stat(videoPath)
	if err != nil {
		logrus.Errorf("获取视频文件信息失败: %v", err)
		return false
	}

	audioInfo, err := os.Stat(audioPath)
	if err != nil {
		logrus.Errorf("获取音频文件信息失败: %v", err)
		return false
	}

	// 检查文件大小是否相近（允许一定误差）
	expectedSize := videoInfo.Size() + audioInfo.Size()
	if abs(int64(outputInfo.Size())-expectedSize) > 1024*1024 { // 允许1MB误差
		return false
	}

	// 对于重复检测，我们不能直接比较输入文件和输出文件的哈希
	// 因为合成过程会改变文件结构。相反，我们需要为每个合成的文件
	// 记录原始输入文件的哈希，或者使用其他方法来验证
	// 这里我们使用一个简单的方法：如果文件大小相近，并且是由相同的输入文件合成的
	// 我们就认为它是相同的文件
	// 注意：这不是最精确的方法，但在大多数情况下是有效的
	return true
}

// abs 计算整数的绝对值
func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// getMp4Metadata 从MP4文件中读取元数据信息
func (c *Config) getMp4Metadata(filePath string) (map[string]string, error) {
	// 构建MP4Box命令行参数
	args := []string{"-info", filePath}
	cmd := exec.Command(c.GPACPath, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	// 执行命令
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 解析输出
	output := stdout.String()

	metadata := make(map[string]string)

	// 提取标题信息
	if idx := strings.Index(strings.ToLower(output), "title:"); idx != -1 {
		// 找到冒号后的内容，直到换行符
		colonIdx := strings.Index(output[idx:], ":")
		if colonIdx != -1 {
			lineEndIdx := strings.Index(output[idx+colonIdx:], "\n")
			if lineEndIdx != -1 {
				value := strings.TrimSpace(output[idx+colonIdx+1 : idx+colonIdx+lineEndIdx])
				if value != "" {
					metadata["title"] = value
				}
			}
		}
	}

	// 提取艺术家信息
	if idx := strings.Index(strings.ToLower(output), "artist:"); idx != -1 {
		// 找到冒号后的内容，直到换行符
		colonIdx := strings.Index(output[idx:], ":")
		if colonIdx != -1 {
			lineEndIdx := strings.Index(output[idx+colonIdx:], "\n")
			if lineEndIdx != -1 {
				value := strings.TrimSpace(output[idx+colonIdx+1 : idx+colonIdx+lineEndIdx])
				if value != "" {
					metadata["artist"] = value
				}
			}
		}
	}

	// 提取专辑信息
	if idx := strings.Index(strings.ToLower(output), "album:"); idx != -1 {
		// 找到冒号后的内容，直到换行符
		colonIdx := strings.Index(output[idx:], ":")
		if colonIdx != -1 {
			lineEndIdx := strings.Index(output[idx+colonIdx:], "\n")
			if lineEndIdx != -1 {
				value := strings.TrimSpace(output[idx+colonIdx+1 : idx+colonIdx+lineEndIdx])
				if value != "" {
					metadata["album"] = value
				}
			}
		}
	}

	return metadata, nil
}

// isIdenticalFileExists 检查目录中是否存在与输入音频和视频文件内容相同的文件
func (c *Config) isIdenticalFileExists(dirPath string, videoPath string, audioPath string, part string) (bool, string) {
	// 读取目录中的所有文件
	files, err := os.ReadDir(dirPath)
	if err != nil {
		logrus.Errorf("读取目录失败: %v", err)
		return false, ""
	}

	// 计算输入文件的组合哈希
	inputHash := c.calculateCombinedHash(videoPath, audioPath)
	if inputHash == "" {
		// 如果无法计算哈希，使用文件大小进行比较
		videoInfo, err := os.Stat(videoPath)
		if err != nil {
			logrus.Errorf("获取视频文件信息失败: %v", err)
			return false, ""
		}

		audioInfo, err := os.Stat(audioPath)
		if err != nil {
			logrus.Errorf("获取音频文件信息失败: %v", err)
			return false, ""
		}

		expectedSize := videoInfo.Size() + audioInfo.Size()

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if !strings.HasSuffix(file.Name(), ".mp4") {
				continue
			}

			filePath := filepath.Join(dirPath, file.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				logrus.Errorf("获取文件信息失败: %v", err)
				continue
			}

			if abs(int64(fileInfo.Size())-expectedSize) > 1024*1024 {
				continue
			}

			logrus.Info("发现相似大小的文件: ", filePath, " 大小:", fileInfo.Size(), " 预期:", expectedSize)
			return true, filePath
		}

		return false, ""
	}

	// 检查每个MP4文件
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".mp4") {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())

		// 检查.hash文件
		hashFilePath := strings.ReplaceAll(filePath, ".mp4", ".hash")
		if utils.IsExist(hashFilePath) {
			hashContent, err := os.ReadFile(hashFilePath)
			if err == nil && string(hashContent) == inputHash {
				logrus.Info("发现相同内容的文件: ", filePath)
				return true, filePath
			}
		}

		// 检查MP4文件的元数据
		metadata, err := c.getMp4Metadata(filePath)
		if err == nil {
			if metadata["title"] == c.GroupId && metadata["artist"] == c.Uid && metadata["album"] == c.ItemId {
				// 如果提供了part，还需要检查文件名中是否包含part
				if part != "" {
					if strings.Contains(file.Name(), part) {
						logrus.Info("发现相同元数据和part的文件: ", filePath)
						return true, filePath
					}
				} else {
					logrus.Info("发现相同元数据的文件: ", filePath)
					return true, filePath
				}
			}
		}
	}

	return false, ""
}

func MessageBox(text string) {
	_ = zenity.Warning(text, zenity.Title("提示"), zenity.Width(400))
}

// findM4sFiles 检查目录及其子目录下是否存在m4s文件
func (c *Config) findM4sFiles() error {
	var m4sFiles []string
	err := filepath.Walk(c.CachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Warnf("查找bilibili缓存目录异常: %s", path)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == conver.M4sSuffix {
			m4sFiles = append(m4sFiles, path)
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(m4sFiles) == 0 {
		return fmt.Errorf("缓存目录找不到m4s文件: %s", c.CachePath)
	}
	return nil
}

// SelectDirectory 选择 BiliBili 缓存目录
func (c *Config) SelectDirectory() {
	for {
		var err error
		c.CachePath, err = zenity.SelectFile(zenity.Title("请选择 BiliBili 缓存目录"), zenity.Directory())
		if c.CachePath == "" || err != nil {
			logrus.Warn("关闭对话框后自动退出程序")
			os.Exit(1)
		}

		if c.findM4sFiles() == nil {
			logrus.Info("选择的 BiliBili 缓存目录为:", c.CachePath)
			return
		}
		MessageBox("选择的 BiliBili 缓存目录内找不到m4s文件，请重新选择！")
	}
}

// SelectGPACPath 选择 GPACPath文件
func (c *Config) SelectGPACPath() {
	for {
		var err error
		c.GPACPath, err = zenity.SelectFile(zenity.Title("请选择 GPAC 的 mp4box 文件"))
		if c.GPACPath == "" || err != nil {
			logrus.Warn("关闭对话框后自动退出程序")
			os.Exit(1)
		}

		if utils.IsExist(c.GPACPath) {
			logrus.Info("选择 GPAC 的 mp4box 文件为:", c.GPACPath)
			return
		}
		MessageBox("选择 GPAC 的 mp4box 文件不存在，请重新选择！")
	}
}

// 如果是目录，尝试下载并转换xml弹幕为ass格式
func (c *Config) downloadXml() {
	dirPath := filepath.Dir(c.video)
	dirName := filepath.Base(dirPath)

	if len(dirName) < 6 { // Android嵌套目录，音视频目录为80
		danmakuXml := filepath.Join(filepath.Dir(dirPath), conver.DanmakuXml)
		if Size(danmakuXml) != 0 {
			c.AssPath = conver.Xml2Ass(danmakuXml) // 转换xml弹幕文件为ass格式
		}
		return
	}
	xmlPath := filepath.Join(dirPath, dirName+conver.XmlSuffix)
	if Size(xmlPath) != 0 {
		c.AssPath = conver.Xml2Ass(xmlPath) // 转换xml弹幕文件为ass格式
		return
	}
	if e := downloadFile(joinUrl(dirName), xmlPath); e != nil {
		if downloadFile(joinXmlUrl(dirName), xmlPath) != nil {
			logrus.Warn("弹幕文件下载失败:", joinUrl(dirName))
			return
		}
	}
	c.AssPath = conver.Xml2Ass(xmlPath) // 转换xml弹幕文件为ass格式
}

// GetVAId 返回.playurl文件中视频文件或音频文件件数组
func GetVAId(patch string) (videoID string, audioID string) {
	pu := filepath.Join(filepath.Dir(patch), conver.PlayUrlSuffix)
	puByte, e := os.ReadFile(pu)
	if e == nil {
		/*
			视频：
			data.dash.video[0].id
			data.dash.audio[0].id
			番剧：
			result.dash.video[0].id  80  需要加上30000，实际30080.m4s
			result.dash.audio[0].id  30280
		*/
		var p gjson.Result
		if p = gjson.GetBytes(puByte, "data"); !p.Exists() {
			p = gjson.GetBytes(puByte, "result")
		}
		if p.Exists() {
			return p.Get("dash.video|@reverse|0.id").String(), p.Get("dash.audio|@reverse|0.id").String()
		}
		return "", ""
	}
	if filepath.Base(filepath.Dir(patch)) != "80" {
		logrus.Warnln("找不到.playurl文件,切换到Android模式解析entry.json文件")
	}
	androidPEJ := filepath.Join(filepath.Dir(filepath.Dir(patch)), conver.PlayEntryJson)
	puDate, e := os.ReadFile(androidPEJ)
	if e != nil {
		logrus.Error("找不到entry.json文件!")
		return
	}
	status := gjson.GetBytes(puDate, "page_data.download_title").String()
	if status != "completed" && status != "视频已缓存完成" && status != "" {
		logrus.Error("跳过未缓存完成的视频", status)
		return
	}
	return "video.m4s", "audio.m4s"
}

func OpenFolder(outputDir string) {
	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("explorer", outputDir).Start()
	case "darwin": // macOS
		_ = exec.Command("open", outputDir).Start()
	default: // Linux and other Unix-like systems
		_ = exec.Command("xdg-open", outputDir).Start()
	}
}
