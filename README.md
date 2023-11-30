# m4s-converter

将bilibili缓存的m4s转成mp4(读PC端缓存目录)

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
"tabName": "正片",
"uname": "珂姬与科技", // 上传的用户名
"title": "蛇的“工作原理”",// 单个视频名称
"groupTitle": "3D动画之工作原理", // 视频组名称
}
```

转换后的文件命名规则：title-uname.mp4

```
文件名识别：
1332097557-1-30280.m4s // 所有30280均为音频文件
1332097557-1-100048.m4s // 值不固定的为视频文件
```

自定义参数:
```
.\m4s-converter.exe -h
Usage of m4s-converter.exe:
  -c string
        指定bilibili缓存目录 (default "C:\\Users\\mzky\\Videos\\bilibili")
  -f string
        指定ffmpeg路径，或将本程序访问ffmpeg.exe同目录 (default "D:\\git\\m4s-converter\\ffmpeg.exe")
  -y    是否覆盖，默认不覆盖
  
.\m4s-converter.exe -f 
.\m4s-converter.exe -f C:\ff\ffmpeg.exe 
.\m4s-converter.exe -f C:\ff\ffmpeg.exe -y
.\m4s-converter.exe -f C:\ff\ffmpeg.exe -c C:\Users\mzky\Videos\bilibili -y
```

验证合成：
```
....略
      handler_name    : SoundHandler
Stream mapping:
  Stream #0:0 -> #0:0 (copy)
  Stream #1:0 -> #0:1 (aac (native) -> aac (native))
Press [q] to stop, [?] for help
Output #0, mp4, to 'C:\Users\mzky\Videos\bilibili\test\2M大小开源神器，Windows 11的救星！-JOKER鹏少.mp4':
  Metadata:
    major_brand     : iso5
    minor_version   : 1
    compatible_brands: avc1iso5dsmsmsixdash
    description     : Packed by Bilibili XCoder v2.0.2
    encoder         : Lavf58.45.100
    Stream #0:0(und): Video: h264 (High) (avc1 / 0x31637661), yuv420p(tv, bt709), 1920x1080 [SAR 1:1 DAR 16:9], q=2-31, 8 kb/s, 30 fps, 30 tbr, 16k tbn, 16k tbc (default)
    Metadata:
      handler_name    : VideoHandler
    Stream #0:1(und): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 128 kb/s (default)
    Metadata:
      handler_name    : SoundHandler
      encoder         : Lavc58.91.100 aac
frame= 6545 fps=2032 q=-1.0 Lsize=    9123kB time=00:03:38.26 bitrate= 342.4kbits/s speed=67.8x
video:5462kB audio:3391kB subtitle:0kB other streams:0kB global headers:0kB muxing overhead: 3.052492%
[aac @ 000002beaa4f9200] Qavg: 754.864
2023/11/30 18:26:37 已合成视频文件： C:\Users\mzky\Videos\bilibili\test\2M大小开源神器，Windows 11的救星！-JOKER鹏少.mp4
2023/11/30 18:26:37 任务已全部完成:
C:\Users\mzky\Videos\bilibili\1332097557\蛇的“工作原理”-珂姬与科技.mp4
C:\Users\mzky\Videos\bilibili\1333045397\2M大小开源神器，Windows 11的救星！-JOKER鹏少.mp4
C:\Users\mzky\Videos\bilibili\1335961247\豆瓣9.1分，时隔40年就再也没人能拍出这么干净的国产电影了！唉！！-7块电影.mp4
C:\Users\mzky\Videos\bilibili\test\2M大小开源神器，Windows 11的救星！-JOKER鹏少.mp4
PS C:\ff>
```
