package common

import (
	"fmt"
	"github.com/bingoohuang/golog"
)

func (c *Config) InitLog() {
	layout := `%t{yyyy-MM-dd_HH:mm:ss} [%-5l{length=5,printColor=true}] %msg{singleLine=false}%n`
	spec := fmt.Sprintf("file=%s,stdout=true", "m4s.log")
	golog.Setup(golog.Layout(layout), golog.Spec(spec))
}
