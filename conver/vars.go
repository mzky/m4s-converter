package conver

var (
	AssSuffix       = ".ass"
	XmlSuffix       = ".xml"
	M4sSuffix       = ".m4s"
	Mp4Suffix       = ".mp4"
	VideoInfoSuffix = ".videoInfo"
	VideoInfoJson   = "videoInfo.json"
	AudioSuffix     = "-audio.mp3"
	VideoSuffix     = "-video.mp4"
	PlayUrlSuffix   = ".playurl"
	/*
			文件名识别：
			1332097557-1-30280.m4s // 所有30280均为音频文件,后来发现还有30216，所以需要从.playurl文件中取
			1332097557-1-100048.m4s // 值不固定的为视频文件
			.playurl文件：
		    data.dash.video[].id
			data.dash.audio[].id
	*/
	AudioFileID = "30280"
)

type PlayUrl struct {
	Data struct {
		Dash struct {
			Video []struct {
				ID     int `json:"id"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"video"`
			Audio []struct {
				ID int `json:"id"`
			} `json:"audio"`
		} `json:"dash"`
	} `json:"data"`
}
