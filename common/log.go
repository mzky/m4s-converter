package common

import (
	"fmt"
	"github.com/bingoohuang/golog"
)

func (c *Config) LogInit() {
	layout := `%t{yyyy-MM-dd_HH:mm:ss.SSS} [%-5l{length=5}] %msg %fields%n`
	spec := fmt.Sprintf("file=%s,stdout=true", "m4s.log")
	golog.Setup(golog.Layout(layout), golog.Spec(spec))
}

func (c *Config) PanicHandler() {
	if e := recover(); e != nil {
		c.File.Close()
		fmt.Print("按回车键退出...")
		fmt.Scanln()
	}
}
