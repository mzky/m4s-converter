## 为什么开发此程序？
bilibili下架了很多视频，之前收藏和缓存的视频均无法播放

![image](https://github.com/mzky/m4s-converter/assets/13345233/ea8bc799-e47d-40ca-bde4-c47193f0e453)

- 喜欢的视频赶紧缓存起来，使用本程序将bilibili缓存的m4s转成mp4，方便随时播放。

- 因bilibili使用的是GPAC处理视频，本工具从v1.5.0开始默认使用GPAC的MP4Box进行音视频合成（此版开始不支持32位系统），能够避免FFMpeg合成视频后音画不同步问题，详见：https://github.com/mzky/m4s-converter/issues/11


### 下载后双击执行或命令行执行(windows、linux版自测正常，MacOS未测试[欢迎反馈])
- https://github.com/mzky/m4s-converter/releases/latest


### Android导出的文件合并方法 
- 详见：https://github.com/mzky/m4s-converter/issues/9


### 除window和linux外，其它环境的依赖工具安装
- https://github.com/mzky/m4s-converter/wiki


### 命令行参数
```
# 指定FFMpeg路径: ./m4s-converter-linux_amd64 -f /var/FFMpeg/ffmpeg 或 ./m4s-converter-amd64 -f select
# 指定MP4Box路径: ./m4s-converter-amd64.exe -g "D:\GPAC\mp4box.exe" 或 ./m4s-converter-amd64 -g select
 Flags: 
    -h --help         查看帮助信息
    -v --version      查看版本信息
    -a --assoff       关闭自动生成弹幕功能，默认不关闭
    -s --skip         跳过合成同名视频(优先级高于overlay)，默认不跳过，但会跳过[完全相同]的文件
    -o --overlay      合成文件时是否覆盖同名视频，默认不覆盖并重命名新文件
    -c --cachepath    自定义视频缓存路径，默认使用bilibili的默认缓存路径
    -g --gpacpath     自定义GPAC的mp4box文件路径,值为select时弹出选择对话框
    -f --ffmpegpath   自定义FFMpeg文件路径,值为select时弹出选择对话框
```


### 验证合成：
```
2023-12-05_16:02:46 [INFO ] 已合成视频文件:中国-美景极致享受-笨蹦崩.mp4
2023-12-05_16:02:46 [INFO ] ==========================================
2023-12-05_16:02:46 [INFO ] 合成的文件:
C:\Users\mzky\Videos\bilibili\output\【获奖学生动画】The Little Poet 小诗人｜CALARTS 2023\【获奖学生动画】The Little Poet 小诗人｜CALARTS 2023-toh糖.mp4
C:\Users\mzky\Videos\bilibili\output\【电影历史_专题片】《影响》致敬中国电影40年【全集】\40年光影记忆-开飞机的巡查司.mp4
C:\Users\mzky\Videos\bilibili\output\“我不是个好导演”，听田壮壮讲述“我和电影的关系”\“我不是个好导演”，听田壮壮讲述“我和电影的关系”-Tatler的朋友们.mp4
C:\Users\mzky\Videos\bilibili\output\【4K8K-世界各地的美景】\中国-美景极致享受-笨蹦崩.mp4
2023-12-05_16:02:46 [INFO ] 已完成本次任务，耗时:5秒
2023-12-05_16:02:46 [INFO ] ==========================================
按回车键退出...
```

- 合成 1.46GB 文件，耗时: 5 秒
- 合成 11.7GB 文件，耗时:38 秒

以上为固态硬盘测试结果, 仅供参考

##
#### 弹幕xml转换为ass使用了此项目
- https://github.com/kafuumi/converter


#### 视频编码使用的工具
- https://gpac.io
- https://ffmpeg.org
- 本程序不会对下载的音视频转码，仅通过上面两个工具进行音视频轨合成


#### 非缓存方式下载，推荐使用其它工具
- https://github.com/nICEnnnnnnnLee/BilibiliDown
- https://github.com/leiurayer/downkyi


## 提缺陷和建议

知乎不常上，缺陷或建议提交 [issues](https://github.com/mzky/m4s-converter/issues/new/choose) , 最好带上异常视频的URL地址
