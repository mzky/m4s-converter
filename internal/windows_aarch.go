//go:build windows && aarch

package internal

var ffMpegName = "ffmpeg.exe"

// 没找到aarch合适的 ffmpeg ，先自己下载放到本地吧
func GetFFMpeg() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, ffMpegName)
}
