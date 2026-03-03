package common

import (
	"compress/flate"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

func downloadFile(url string, filepath string) error {
	// 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: 3 * time.Second, // 3秒超时
	}
	// 发起HTTP GET请求
	httpReq, err := client.Get(url)
	if err != nil {
		return errors.Wrap(err, "HTTP请求失败")
	}
	defer httpReq.Body.Close()

	if httpReq.StatusCode != http.StatusOK {
		return errors.New("无法获取字幕数据")
	}

	// 创建本地文件
	localFile, err := os.Create(filepath)
	if err != nil {
		return err
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
			return errors.New("无法获取字幕数据")
		}

		// 将解压后的数据写入本地文件
		if _, err := localFile.Write(bodyBytes); err != nil {
			return err
		}
	} else {
		// 如果不是deflate编码，直接将响应体写入文件
		if _, err := io.Copy(localFile, httpReq.Body); err != nil {
			return err
		}
	}

	// 检查文件是否成功写入
	if e := localFile.Sync(); e != nil {
		return e
	}

	return nil
}
