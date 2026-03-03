package common

import (
	"compress/flate"
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

func downloadFile(ctx context.Context, url string, filepath string) error {
	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 使用 context 创建请求，支持取消和超时
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "创建 HTTP 请求失败")
	}

	// 发起HTTP GET请求
	httpReq, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "下载字幕文件失败")
	}
	defer httpReq.Body.Close()

	if httpReq.StatusCode != http.StatusOK {
		return errors.Errorf("无法获取字幕数据，HTTP 状态码: %d", httpReq.StatusCode)
	}

	// 创建本地文件
	localFile, err := os.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "创建本地文件失败")
	}
	defer localFile.Close()

	// 检查Content-Encoding是否为deflate
	contentEncoding := httpReq.Header.Get("Content-Encoding")
	if contentEncoding == "deflate" {
		// 如果是deflate编码，解压缩数据
		reader := flate.NewReader(httpReq.Body)
		defer reader.Close()

		// 读取并解压数据
		bodyBytes, err := io.ReadAll(reader)
		if err != nil || bodyBytes == nil {
			return errors.New("解压缩字幕数据失败")
		}

		// 将解压后的数据写入本地文件
		if _, err := localFile.Write(bodyBytes); err != nil {
			return errors.Wrap(err, "写入字幕文件失败")
		}
	} else {
		// 如果不是deflate编码，直接将响应体写入文件
		if _, err := io.Copy(localFile, httpReq.Body); err != nil {
			return errors.Wrap(err, "写入字幕文件失败")
		}
	}

	// 检查文件是否成功写入
	if e := localFile.Sync(); e != nil {
		return errors.Wrap(e, "同步字幕文件失败")
	}

	return nil
}
