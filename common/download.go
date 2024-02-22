package common

import (
	"compress/flate"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, filepath string) error {
	// 发起HTTP GET请求
	httpReq, err := http.Get(url)
	if err != nil {
		return err
	}
	defer httpReq.Body.Close()

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
		if err != nil {
			return err
		}

		// 将解压后的数据写入本地文件
		if _, err := localFile.Write(bodyBytes); err != nil {
			return err
		}
	} else {
		// 如果不是deflate编码，直接将Content-Encoding写入文件
		if _, err := localFile.Write([]byte(contentEncoding)); err != nil {
			return err
		}
	}

	// 检查文件是否成功写入
	if err := localFile.Sync(); err != nil {
		return err
	}

	return nil
}
