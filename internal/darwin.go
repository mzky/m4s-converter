//go:build darwin

package internal

func GetMP4Box() string {
	return getCliPath("MP4Box")
}
