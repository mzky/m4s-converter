package common

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/integrii/flaggy"
	"os"
	"strings"
)

var (
	Subtitle *flaggy.Subcommand
)

func appendExp(exp *[]string, format string, a ...any) {
	*exp = append(*exp, fmt.Sprintf(format, a...))
}

func subtitleDescription() string {
	var exp []string
	appendExp(&exp, "subtitle 自动下载xml弹幕文件并转换为ass")
	appendExp(&exp, "\t\t   %s sub", os.Args[0])
	return strings.Join(exp, "\n")
}

func InitFlags() {
	flaggy.DefaultParser.DisableShowVersionWithVersion()
	flaggy.SetName(color.CyanString("m4s-converter"))
	flaggy.SetDescription(color.CyanString("Bilibili Tool"))

	Subtitle = flaggy.NewSubcommand("sub")
	flaggy.AttachSubcommand(Subtitle, 1)
	Subtitle.Description = subtitleDescription()

	// 解析到未预期参数时显示帮助
	flaggy.ShowHelpOnUnexpectedEnable()
	flaggy.Parse()
}
