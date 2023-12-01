package common

import (
	"fmt"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

func _TEXT(str string) *uint16 {
	ptr, _ := syscall.UTF16PtrFromString(str)
	return ptr
}

func (c *Config) MessageBox(text string) {
	var handle windows.Handle
	log.Println(text)
	win.MessageBox(win.HWND(handle), _TEXT(text), _TEXT("消息"), win.MB_ICONWARNING)
}

func (c *Config) SelectFile() {
	fileName := make([]uint16, 256)
	ofn := win.OPENFILENAME{
		LStructSize: uint32(unsafe.Sizeof(win.OPENFILENAME{})),
		LpstrFile:   &fileName[0],
		NMaxFile:    uint32(len(fileName)),
		HwndOwner:   win.HWND(0),
		//LpstrFilter:   _TEXT("选择文件 (*.exe)\000*.exe\000"),
		LpstrTitle:    _TEXT("选择ffmpeg.exe文件"),
		NMaxFileTitle: 255,
		//Flags:         win.OFN_EXPLORER | win.OFN_FILEMUSTEXIST | win.OFN_PATHMUSTEXIST | win.OFN_LONGNAMES,
	}

	// 打开文件选择窗口
	if !win.GetOpenFileName(&ofn) {
		fmt.Println("关闭对话框后自动退出程序")
		os.Exit(1)
	}

	c.FFmpegPath = win.UTF16PtrToString(ofn.LpstrFile)
	if strings.Contains(c.FFmpegPath, "ffmpeg.exe") {
		log.Println("选择的ffmpeg文件: ", c.FFmpegPath)
		return
	}
	c.MessageBox("请选择 ffmpeg.exe 程序文件，如未安装可从以下地址下载：" +
		"\n https://github.com/GyanD/codexffmpeg/releases" +
		"\n https://github.com/BtbN/FFmpeg-Builds/releases")
	c.SelectFile() // 选错重新弹出对话框进行选择
}

func (c *Config) SelectDirectory() {
	var bsi win.BROWSEINFO
	bsi.LpszTitle = _TEXT("请选择 bilibili 缓存目录")

	pid := win.SHBrowseForFolder(&bsi)
	if pid == 0 {
		fmt.Println("关闭对话框后自动退出程序")
		os.Exit(1)
	}
	defer win.CoTaskMemFree(pid)

	path := make([]uint16, win.MAX_PATH)
	win.SHGetPathFromIDList(pid, &path[0])

	c.CachePath = syscall.UTF16ToString(path)
	if Exist(filepath.Join(c.CachePath, ".videoInfo")) || Exist(filepath.Join(c.CachePath, "load_log")) {
		log.Println("选择的 bilibili 缓存目录为:", c.CachePath)
		return
	}
	c.MessageBox("选择的 bilibili 缓存目录不正确，请重新选择！")
	c.SelectDirectory()
}
