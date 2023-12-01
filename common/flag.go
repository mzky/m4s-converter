package common

import "flag"

func (c *Config) Flags() {
	ffmpegPath := flag.String("f", c.FFmpegPath, "指定ffmpeg文件路径")
	cachePath := flag.String("c", c.CachePath, "指定bilibili缓存目录")
	overlay := flag.Bool("y", false, "是否覆盖，默认不覆盖")

	flag.Parse()

	if *overlay {
		c.Overlay = "-y"
	}

	if *ffmpegPath != c.FFmpegPath {
		c.FFmpegPath = *ffmpegPath
	}

	if *cachePath != c.CachePath {
		c.CachePath = *cachePath
	}
}
