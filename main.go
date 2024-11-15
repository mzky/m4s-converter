package main

import (
	"m4s-converter/common"
)

func main() {
	var c common.Config
	go c.PanicHandler()
	c.InitLog()
	c.InitConfig()
	c.Synthesis()
}
