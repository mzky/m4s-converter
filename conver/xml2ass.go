package conver

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"

	"github.com/mzky/converter"
)

func Xml2ass(xml string) string {
	var dstFile string
	xmlState, err := os.Stat(xml)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Fatalf("文件：%s不存在\n", xml)
		} else {
			logrus.Fatalln(err)
		}
	}

	xmls := make([]string, 0)
	if xmlState.IsDir() {
		if xml[len(xml)-1] != os.PathSeparator {
			xml += string(os.PathSeparator)
		}
		if entries, e := os.ReadDir(xml); e != nil {
			logrus.Error(e)
		} else {
			for _, entry := range entries {
				if !entry.IsDir() {
					name := entry.Name()
					if strings.HasSuffix(name, XmlSuffix) {
						xmls = append(xmls, xml+name)
					}
				}
			}
		}
	} else {
		if strings.HasSuffix(xml, XmlSuffix) {
			xmls = append(xmls, xml)
		} else {
			logrus.Fatalln("不支持的文件格式。")
		}
	}

	setting := DefaultSetting
	assConfig := setting.GetAssConfig()
	chain := converter.NewFilterChain()
	keywordFilter, typeFilter := setting.GetFilter()
	chain.AddFilter(keywordFilter).AddFilter(typeFilter)
	var success int32 = 0
	var failed int32 = 0
	for _, file := range xmls {
		//加载xml文件
		src, _ := os.Open(file)
		if src == nil {
			failed++
			return dstFile
		}
		//如果在go程中加载xml，当文件过多时会出现过高的内存占用
		pool := converter.LoadPool(src, chain)
		_ = src.Close()

		dstFile = strings.ReplaceAll(file, filepath.Ext(file), AssSuffix)
		dst, e := os.Create(dstFile)
		if e != nil {
			failed++
			return dstFile
		}
		if er := pool.Convert(dst, assConfig); er == nil {
			//fmt.Printf("[ok] %s ==> %s\n", file, dstFile)
			success++
		} else {
			failed++
			//fmt.Printf("[failed] %s\n", file)
		}
		_ = dst.Close()
	}
	fmt.Println("转换弹幕:", "成功数", success, "失败数", failed)
	return dstFile
}
