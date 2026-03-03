package main

import (
	"m4s-converter/common"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	var c common.Config
	c.InitLog()
	c.InitConfig()

	// 捕获 SIGINT 信号（Ctrl+C）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// 在 goroutine 中等待信号
	go func() {
		<-sigChan
		logrus.Info("收到退出信号，正在处理当前任务...")
		c.SetExitFlag(true)
	}()

	c.Synthesis()
}
