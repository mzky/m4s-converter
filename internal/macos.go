//go:build macos

package internal

func GetFFMpeg() string {
	return GetCliPath("ffmpeg")
}
