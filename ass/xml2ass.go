package ass

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
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
		if entries, err := os.ReadDir(xml); err != nil {
			logrus.Fatalln(err)
		} else {
			for _, entry := range entries {
				if !entry.IsDir() {
					name := entry.Name()
					if strings.HasSuffix(name, ".xml") {
						xmls = append(xmls, xml+name)
					}
				}
			}
		}
	} else {
		if strings.HasSuffix(xml, ".xml") {
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
		dotIndex := strings.LastIndex(file, ".")
		if dotIndex == -1 {
			dotIndex = len(file)
		}
		dstFile = file[:dotIndex] + ".ass"
		dst, err := os.Create(dstFile)
		if err != nil {
			failed++
			logrus.Println(err)
			return dstFile
		}
		if err := pool.Convert(dst, assConfig); err == nil {
			//fmt.Printf("[ok] %s ==> %s\n", file, dstFile)
			success++
		} else {
			failed++
			//fmt.Printf("[failed] %s\n", file)
		}
	}
	fmt.Printf("转换ass弹幕：%d, 转换成功数：%d 转换失败数：%d\n", len(xmls), success, failed)
	return dstFile
}
