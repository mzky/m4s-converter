# m4s-converter

将bilibili缓存的m4s转成mp4(读PC端缓存目录)

```
支持单独和批量合成，比如：
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

转换后的文件命名规则：groupTitle-title-uname.mp4

文件名识别：
1332097557-1-30280.m4s // 所有30280均为音频
1332097557-1-100048.m4s // 值不固定

