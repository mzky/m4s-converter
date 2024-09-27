//go:build linux

package internal

func GetFFMpeg() string {
	return GetCliPath("ffmpeg")
}
