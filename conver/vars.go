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
			.playurl文件：
		    data.dash.video[].id
			data.dash.audio[].id
	*/
	// AudioFileID = "30280"
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

type Entry struct {
	Title     string `json:"title"`
	TypeTag   string `json:"type_tag"`
	Avid      int64  `json:"avid"`
	Bvid      string `json:"bvid"`
	OwnerId   int    `json:"owner_id"`
	OwnerName string `json:"owner_name"`
	PageData  struct {
		Cid              int64  `json:"cid"`
		Part             string `json:"part"`
		DownloadTitle    string `json:"download_title"`    // "视频已缓存完成"
		DownloadSubtitle string `json:"download_subtitle"` // "牢中注定（1955）"
	} `json:"page_data"`
}
