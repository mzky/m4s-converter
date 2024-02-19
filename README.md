# m4s-converter

## 为什么开发此程序？
bilibili下架了很多视频，之前收藏和缓存（ipad）的均无法播放，

喜欢的视频赶紧缓存起来，使用本程序将bilibili缓存的m4s转成mp4，以便后续播放

```
支持单独和批量目录识别，比如：
C:\Users\mzky\Videos\bilibili\1332097557
C:\Users\mzky\Videos\bilibili\
```

```
{
"groupId": 1582186,
"itemId ":1332097557,
"aid": 620879911,
"cid": 1332097557, // 文件夹名
"uname": "珂姬与科技", // 上传的用户名
"title": "蛇的“工作原理”",// 单个视频名称
"groupTitle": "3D动画之工作原理", // 视频组名称
"status": "downloading",// downloading正在缓存中；pending等待缓存（还没有缓存文件）；completed缓存完成
}
```

转换后的文件命名规则：title-uname.mp4

```
文件名识别：
1332097557-1-30280.m4s // 所有30280均为音频文件
1332097557-1-100048.m4s // 值不固定的为视频文件
```

验证合成：
```
2023-12-05_16:02:41 [INFO ] 已将m4s转换为音视频文件:C:\Users\mzky\Videos\bilibili\1120254313\1120254313_nb3-1-30080.m4s-video.mp4
2023-12-05_16:02:41 [INFO ] 已将m4s转换为音视频文件:C:\Users\mzky\Videos\bilibili\1120254313\1120254313_nb3-1-30280.m4s-audio.mp3
2023-12-05_16:02:41 [INFO ] 已将m4s转换为音视频文件:C:\Users\mzky\Videos\bilibili\65093887\65093887-1-30032.m4s-video.mp4
2023-12-05_16:02:41 [INFO ] 已将m4s转换为音视频文件:C:\Users\mzky\Videos\bilibili\65093887\65093887-1-30280.m4s-audio.mp3
2023-12-05_16:02:42 [INFO ] 已将m4s转换为音视频文件:C:\Users\mzky\Videos\bilibili\799281779\799281779_nb3-1-30080.m4s-video.mp4
2023-12-05_16:02:42 [INFO ] 已将m4s转换为音视频文件:C:\Users\mzky\Videos\bilibili\799281779\799281779_nb3-1-30280.m4s-audio.mp3
2023-12-05_16:02:43 [INFO ] 已将m4s转换为音视频文件:C:\Users\mzky\Videos\bilibili\869752798\869752798_da2-1-30080.m4s-video.mp4
2023-12-05_16:02:43 [INFO ] 已将m4s转换为音视频文件:C:\Users\mzky\Videos\bilibili\869752798\869752798_da2-1-30280.m4s-audio.mp3
准备合成mp4 .............
2023-12-05_16:02:43 [INFO ] 已合成视频文件:【获奖学生动画】The Little Poet 小诗人｜CALARTS 2023-toh糖.mp4
准备合成mp4 ................
2023-12-05_16:02:43 [INFO ] 已合成视频文件:40年光影记忆-开飞机的巡查司.mp4
准备合成mp4 ................
2023-12-05_16:02:45 [INFO ] 已合成视频文件:“我不是个好导演”，听田壮壮讲述“我和电影的关系”-Tatler的朋友们.mp4
准备合成mp4 ...............
2023-12-05_16:02:46 [INFO ] 已合成视频文件:中国-美景极致享受-笨蹦崩.mp4
2023-12-05_16:02:46 [INFO ] ==========================================
2023-12-05_16:02:46 [INFO ] 合成的文件:
C:\Users\mzky\Videos\bilibili\output\【获奖学生动画】The Little Poet 小诗人｜CALARTS 2023\【获奖学生动画】The Little Poet 小诗人｜CALARTS 2023-toh糖.mp4
C:\Users\mzky\Videos\bilibili\output\【电影历史_专题片】《影响》致敬中国电影40年【全集】\40年光影记忆-开飞机的巡查司.mp4
C:\Users\mzky\Videos\bilibili\output\“我不是个好导演”，听田壮壮讲述“我和电影的关系”\“我不是个好导演”，听田壮壮讲述“我和电影的关系”-Tatler的朋友们.mp4
C:\Users\mzky\Videos\bilibili\output\【4K 8K- 世界各地的美景】\中国-美景极致享受-笨蹦崩.mp4
2023-12-05_16:02:46 [INFO ] 已完成本次任务，耗时:5秒
2023-12-05_16:02:46 [INFO ] ==========================================
按回车键退出...
```

合成 1.46GB 文件，耗时: 5 秒

合成 11.7GB 文件，耗时:38 秒

以上为固态硬盘测试结果

### 非缓存下载方式，推荐使用其它工具
```
https://github.com/nICEnnnnnnnLee/BilibiliDown
https://github.com/leiurayer/downkyi
```

## 弹幕xml转换为ass使用此项目： https://github.com/kafuumi/converter