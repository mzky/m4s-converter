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
		logrus.Errorf("cmd run error, err=%v stderr=%v", err, stdout.String())
		return stdout.String()
	}

	return strings.TrimSpace(stdout.String())
}

func exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
