package common

import (
	"compress/flate"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, filepath string) error {
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	// 创建文件
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 检查Content-Encoding是否为deflate
	contentEncoding := req.Header.Get("Content-Encoding")
	if contentEncoding == "deflate" {
		reader := flate.NewReader(req.Body)
		defer reader.Close()
		// 读取并解压数据
		bodyBytes, e := io.ReadAll(reader)
		if e != nil {
			return e
		}
		file.Write(bodyBytes)
	} else {
		file.Write([]byte(contentEncoding))
	}

	return nil
}
