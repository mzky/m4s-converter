package main

import (
	"context"
	"m4s-converter/common"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建可响应系统信号的 context，支持优雅关闭
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,    // Ctrl+C
		syscall.SIGTERM, // kill 命令
	)
	defer stop()

	var c common.Config
	c.InitLog()
	c.InitConfig(ctx)

	// 使用 context 运行合成，支持取消操作
	if err := c.Synthesis(ctx); err != nil {
		// 程序已处理错误或用户取消，无需额外操作
		return
	}
}
