//go:build darwin

package internal

func GetFFMpeg() string {
	return GetCliPath("ffmpeg")
}
