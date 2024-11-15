### 为什么开发此程序？
bilibili下架了很多视频，之前收藏和缓存的视频均无法播放

![image](https://github.com/mzky/m4s-converter/assets/13345233/ea8bc799-e47d-40ca-bde4-c47193f0e453)

喜欢的视频赶紧缓存起来，使用本程序将bilibili缓存的m4s转成mp4，方便随时播放

### 下载使用(windows、linux版自测正常，MacOS未测试)
https://github.com/mzky/m4s-converter/releases/latest

下载后可直接执行，双击或命令行运行即可
```
程序在windows下能够自动识别默认的bilibli缓存目录，比如：
C:\Users\mzky\Videos\bilibili\
其它系统或者自定义的bilibili缓存路径，请根据提示手动选择目录
```
### 支持Android导出的文件合并
详见：https://github.com/mzky/m4s-converter/issues/9

### 个别系统需要手动安装ffmpeg（桌面版linux系统默认已安装），或指定ffmpeg路径
```
# UOS/Kylin/Ubuntu/Debian等桌面版系统
sudo apt-get install ffmpeg

# OpenEuler/CentOS8等
yum install ffmpeg

# Mac OS
brew install ffmpeg
```
#### 完整版ffmpeg下载地址：https://github.com/BtbN/FFmpeg-Builds/releases
#### 使用GPAC替代ffmpeg合成文件：https://gpac.io/downloads/gpac-nightly-builds
bilibili使用了开源工具GPAC，合成音视频更适合使用GPAC，但GPAC需要下载安装，所以默认使用ffmpeg，
如遇合成的视频音画不同步时，可以使用 -g 参数选择mp4box代替ffmpeg合成文件，详见：https://github.com/mzky/m4s-converter/issues/11

```
# 指定FFMpeg路径: ./m4s-converter-amd64 -f /var/FFMpeg/ffmpeg
# 使用GPAC替代ffmpeg: ./m4s-converter-amd64 -g "C:\Program Files (x86)\GPAC\mp4box.exe" 或 ./m4s-converter-amd64 -g select

Options:
  --assOFF, -a
        是否关闭自动生成ass弹幕，默认不关闭
  --cachePath, -c string
        自定义缓存路径，默认使用BiliBili的默认路径 (default C:\Users\mzky\Videos\bilibili)
  --ffMpeg, -f string
        自定义FFMpeg文件路径
  --gpac, -g string
        使用GPAC的mp4box文件，替代FFMpeg合成文件，输入select则弹出对话框选择文件
  --help, -h
        帮助信息
  --overlay, -o
        是否覆盖已存在的视频，默认不覆盖
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

#### 非缓存方式下载，推荐使用其它工具
```
https://github.com/nICEnnnnnnnLee/BilibiliDown
https://github.com/leiurayer/downkyi
```

#### 弹幕xml转换为ass使用了此项目
```
https://github.com/kafuumi/converter
```

#### 视频编码工具使用了ffmpeg
```
https://ffmpeg.org/
```

### 提缺陷和建议

知乎不常上，缺陷或建议提交 [issues](https://github.com/mzky/m4s-converter/issues/new/choose) 

最好带上异常视频的URL地址
