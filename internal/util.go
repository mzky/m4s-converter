package internal

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// getCliPath 获取命令绝对路径
func getCliPath(name string) string {
	execCmd := exec.Command("which", name)

	var stdout bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stdout

	if err := execCmd.Run(); err != nil {
		logrus.Errorf("找不到MP4Box命令,安装GPAC后重试: %v, %v", err, stdout.String())
		os.Exit(1)
	}

	return strings.TrimSpace(stdout.String())
}
