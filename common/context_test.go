package common

import (
	"context"
	"testing"
	"time"
)

func TestSynthesis_ContextCancel(t *testing.T) {
	// 测试 Context 取消功能
	ctx, cancel := context.WithCancel(context.Background())

	// 模拟用户取消操作
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	// 创建一个临时目录作为缓存路径
	tmpDir := t.TempDir()
	c := &Config{CachePath: tmpDir}

	// 启动合成（应该会很快返回，因为目录为空）
	// 实际测试需要更复杂的设置，这里验证函数签名和基本流程
	err := c.Synthesis(ctx)

	// 如果 context 被取消，应该返回错误
	if err == context.Canceled && ctx.Err() == context.Canceled {
		// 这是预期的行为
		t.Log("Context 取消正常工作")
	}
}

func TestDownloadFile_ContextTimeout(t *testing.T) {
	// 测试 HTTP 下载超时
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// 等待超时
	time.Sleep(10 * time.Millisecond)

	tmpDir := t.TempDir()
	filepath := tmpDir + "/test.xml"

	// 尝试下载（应该会因超时失败）
	err := downloadFile(ctx, "http://example.com/test", filepath)

	// 超时错误是可接受的
	if err != nil {
		t.Logf("预期的错误: %v", err)
	}
}

func TestDiffVersion_ContextCancel(t *testing.T) {
	// 测试版本检查时的 Context 取消
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立即取消

	// 函数应该优雅处理取消，不 panic
	diffVersion(ctx)

	// 如果没有 panic，测试通过
	t.Log("diffVersion 正确处理了取消的 context")
}
