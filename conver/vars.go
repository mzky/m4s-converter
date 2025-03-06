package conver

var (
	AssSuffix         = ".ass"
	XmlSuffix         = ".xml"
	M4sSuffix         = ".m4s"
	Mp4Suffix         = ".mp4"
	VideoInfoSuffix   = ".videoInfo"
	VideoInfoJson     = "videoInfo.json"
	AudioSuffix       = "-audio.mp3"
	AudioSuffixOffset = "_offset-audio.mp3"
	AudioSuffixTS     = "-audio.ts"
	VideoSuffix       = "-video.mp4"
	PlayUrlSuffix     = ".playurl"
	PlayEntryJson     = "entry.json"  // 安卓手机端文件信息
	DanmakuXml        = "danmaku.xml" // 安卓手机端字幕
	/*
			文件名识别：
			1332097557-1-30280.m4s // 所有30280均为音频文件,后来发现还有30216，所以需要从.playurl文件中取
			1332097557-1-100048.m4s // 值不固定的为视频文件

		    视频：
			data.dash.video[0].id
			data.dash.audio[0].id

			番剧：
			result.dash.video[0].id  80  需要加上30000，实际30080.m4s
			result.dash.audio[0].id  30280
	*/
	// AudioFileID = "30280"
)
