# pupucar
放屁车是使用帧同步方案制作的多人在线赛车类游戏

## 关键词
- pupucar
- multiplayer online
- frame alignment
- console game
- byebyebruce/lockstepserver
- hajimehoshi/ebiten


## 内容介绍
1.example.exe为帧同步方案服务端
2.pupucar.exe为放屁车客户端

## 快速启动(windows平台)
1.下载releases文件包

2.启动服务端
```
example.exe
```
3.创建房间

浏览器访问http://localhost:10002/ 点击click按钮创建房间。

4.启动玩家1
```
pupucar.exe -mid 1 -ip 192.168.16.152 
```
其中ip为帧同步服务器IP

5.启动玩家2
```
pupucar.exe -mid 1 -ip 192.168.16.152
```

注意：创建房间及客户端启动应保持在几秒内完成，否则超时房间将自动关闭