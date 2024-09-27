package internal

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

// GetCliPath 获取命令绝对路径，CentOS和Ubuntu地址不同
func GetCliPath(name string) string {
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
