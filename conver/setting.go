package conver

import (
	"encoding/json"
	"github.com/mzky/converter"
	"github.com/sirupsen/logrus"
	"io"
)

// DefaultSetting 默认设置
var DefaultSetting = Setting{
	Fontsize:     26,
	FontName:     "黑体",
	Alpha:        0.3,
	OutlineColor: color{RGB: "0x49516A", Alpha: 0.1},
	ShadowColor:  color{RGB: "0x49516A", Alpha: 0.1},
	RollTime:     15,
	FixTime:      5,
	TimeShift:    0,
	Bold:         true,
	Outline:      0,
	Shadow:       1,
	Width:        1920,
	Height:       1080,
	RollRange:    1.0,
	FixedRange:   1.0,
	Spacing:      0,
	Density:      0,
	Overlay:      false,
	Keyword:      nil,
	Convert:      "s -> r",
}

type color struct {
	RGB   string  `json:"rgb"`
	Alpha float32 `json:"alpha"`
}
type Setting struct {
	Fontsize     int      `json:"fontsize"`     // 字体大小
	FontName     string   `json:"fontName"`     // 字体名称
	Alpha        float32  `json:"alpha"`        // 弹幕透明度
	OutlineColor color    `json:"outlineColor"` // 弹幕描边颜色
	ShadowColor  color    `json:"shadowColor"`  // 弹幕阴影颜色
	RollTime     int      `json:"rollTime"`     // 滚动弹幕显示时间
	FixTime      int      `json:"fixTime"`      // 顶部弹幕和底部弹幕显示时间
	TimeShift    int      `json:"timeShift"`    // 时间偏移,单位秒
	Bold         bool     `json:"bold"`         // 是否粗体
	Outline      int      `json:"outline"`      // 描边大小
	Shadow       int      `json:"shadow"`       // 阴影大小
	Width        int      `json:"width"`        // 视频分辨率的宽
	Height       int      `json:"height"`       // 视频分辨率的高
	RollRange    float32  `json:"rollRange"`    // 滚动弹幕显示范围
	FixedRange   float32  `json:"fixedRange"`   // 顶部弹幕和底部弹幕的显示范围
	Spacing      int      `json:"spacing"`      // 弹幕的上下间距
	Density      int      `json:"density"`      // 同屏弹幕密度,0为无限密度
	Overlay      bool     `json:"overlay"`      // 是否允许弹幕重叠
	Keyword      []string `json:"keyword"`      // 按关键字屏蔽
	Convert      string   `json:"convert"`      // 转换弹幕类型
}

func (s Setting) GetAssConfig() converter.AssConfig {
	textColor, _ := converter.ParseStringARGB(s.Alpha, "0xffffff")
	outlineColor, err := converter.ParseStringARGB(s.OutlineColor.Alpha, s.OutlineColor.RGB)
	if err != nil {
		logrus.Printf("描边颜色设置错误：%v\n", err)
		outlineColor = 0x1E49516A
	}
	shadowColor, err := converter.ParseStringARGB(s.ShadowColor.Alpha, s.ShadowColor.RGB)
	if err != nil {
		logrus.Printf("阴影颜色设置错误：%v\n", err)
		shadowColor = 0x1E49516A
	}

	config := converter.AssConfig{
		Fontsize:     s.Fontsize,
		FontName:     s.FontName,
		Color:        textColor,
		OutlineColor: outlineColor,
		BackColor:    shadowColor,
		RollSpeed:    s.RollTime,
		FixTime:      s.FixTime,
		TimeShift:    s.TimeShift,
		IsBold:       s.Bold,
		Outline:      s.Outline,
		Shadow:       s.Shadow,
		Width:        s.Width,
		Height:       s.Height,
		RollRange:    s.RollRange,
		FixedRange:   s.FixedRange,
		Spacing:      s.Spacing,
		Density:      s.Density,
		Overlay:      s.Overlay,
	}
	return config
}

func (s Setting) GetFilter() (keyword converter.BulletChatFilter,
	convert converter.BulletChatFilter) {
	if s.Keyword == nil {
		keyword = nil
	} else {
		keyword = &converter.KeyWordFilter{Keyword: s.Keyword}
	}
	if s.Convert == "" {
		convert = nil
	} else {
		convert = converter.NewTypeConverter(s.Convert)
	}
	return
}

func ReadSetting(src io.Reader) Setting {
	setting := DefaultSetting
	err := json.NewDecoder(src).Decode(&setting)
	if err != nil && err != io.EOF {
		logrus.Fatalln(err)
	}
	return setting
}
