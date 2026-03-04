## 为什么开发此程序？
bilibili下架了很多视频，之前收藏和缓存的视频均无法播放

![image](https://github.com/mzky/m4s-converter/assets/13345233/ea8bc799-e47d-40ca-bde4-c47193f0e453)

- 收藏的视频建议缓存起来，使用本程序将缓存的音视频m4s文件合并成mp4，方便再次播放。

- 本工具使用GPAC的MP4Box进行音视频合成(v1.5.0之前版本支持FFMpeg)，能够确保合成后的视频质量和同步性。


### 下载后双击执行或通过命令行执行，需要可执行权限
- https://github.com/mzky/m4s-converter/releases/latest


### Android手机端合并文件方法 
- 详见：[拷贝文件与合成方法](https://github.com/mzky/m4s-converter/issues/9)


### 除window和linux外，其它环境的依赖工具安装
- 详见：[依赖工具安装](https://github.com/mzky/m4s-converter/wiki/%E4%BE%9D%E8%B5%96%E5%B7%A5%E5%85%B7%E5%AE%89%E8%A3%85)


### 命令行参数
```
# 指定MP4Box路径: ./m4s-converter-amd64.exe -g "D:\GPAC\mp4box.exe" 或 ./m4s-converter-amd64 -g select
 Flags: 
    -h --help         查看帮助信息
    -v --version      查看版本信息
    -a --assoff       关闭自动生成弹幕功能，默认不关闭
    -o --overlay      合成文件时是否覆盖同名视频，默认不覆盖并重命名新文件
    -u --summarize    将未合并的MP3和视频文件放入汇总目录，默认不汇总
    -c --cachepath    自定义视频缓存路径，默认使用bilibili的默认缓存路径
    -g --gpacpath     自定义GPAC的mp4box文件路径,值为select时弹出选择对话框
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


#### 视频合成使用的工具
- https://gpac.io
- 本程序使用GPAC的MP4Box进行音视频合成，不会对下载的音视频进行转码


#### 本工具无下载视频功能


## 提缺陷和建议

缺陷或建议提交 [issues](https://github.com/mzky/m4s-converter/issues/new/choose) , 最好带上异常视频的URL地址

## Star History

<a href="https://www.star-history.com/#mzky/m4s-converter&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=mzky/m4s-converter&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=mzky/m4s-converter&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=mzky/m4s-converter&type=date&legend=top-left" />
 </picture>
</a>


## ⚠️ **法律声明**  
使用本工具即表示您同意[免责声明](免责声明.md)。  
仅允许转换您本人在视频下架前通过官方客户端合法缓存的内容，且转换结果**严格限于个人备份**，禁止传播与商用。



