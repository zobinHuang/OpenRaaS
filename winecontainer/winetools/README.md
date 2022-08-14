# winecontainer 接口工具

## 介绍

winecounter 核心功能组件，用于 wine 容器的配置与启动

## 软件架构

使用 dockerfile 创建镜像，通过为 `sh` 脚本传参来指定应用

## 安装教程

1. ubuntu 用户可以直接使用 `env_installer_ubuntu.sh` 脚本进行快速环境部署，如果本地有 docker 环境不建议直接运行该脚本
2. 其他用户请参考该文件自行安装对应版本的软件

## 使用说明

1. 创建镜像：`sudo sh build-wine.sh`
2. 将存放应用目录的仓库挂载到 `./apps`
    a.  挂载方法示例：`sudo expect ./auto-mount.exp 1 davfs ip:port /games username pwd`
    b.  取消挂载：`sudo umount ./apps`
3. 运行容器：`sudo sh run-wine.sh image_name container_id apppath appfile appname hwkey screenwidth screenheight targethost wineoptions`
    1. image_name: dcwine 的容器名称，通过设置 `ip:port/dcwine:tag` 指定不同的 depositary
    2. container_id: 容器编号
    3. apppath: 应用目录在 `./apps` 目录下的相对路径
    4. appfile: 应用目录中的可执行文件名
    5. appname: 应用名称，将作为窗口名进行展示
    6. hwkey: 应用类型，传入 "game" 或 "app"
    7. screenwidth: 屏幕宽度
    8. screenheight: 屏幕高度
    9. targethost: RTP 目的主机 IP
    10. fps: 输出帧数
    11. vcodec: 解码方式 h264/vpx
    12. wineoptions: windows app 的可选参数

## 接口文档

| 端口号 | 协议 | 监听者 | 功能 | 文件 |
|  :----:  |  :----:  |  :----:  |  :----:  | :----:  |
| 1xx09 | tcp | 后台服务器 | docker 接收来自后台的按键输入 | syncinput.cpp |
| 1xx05 | rtp | 后台服务器 | docker 向后台推视频流 | supervisord.conf |
| 1xx01 | rtp | 后台服务器 | docker 向后台推音频流 | supervisord.conf |

xx 表示容器编号，如：appvm1->01 appvm23->23

| 内容 | 格式 | 详情 |
|  :----:  |  :----:  |  :----:  |
| 键盘 | "K{KeyCode: int},{KeyState: bool}\|" | KeyCode 为[虚拟键值表编码](https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes), KeyState 为按键状态 (0-弹起, 1-按下) |
| 鼠标 | "M{IsLeft: byte},{mouseState: byte},{X: float},{Y: float},{Width: float},{Height: float}\|" | IsLeft 判断是否左键的操作 (0-否, 非0-是), mouseState 为鼠标操作状态 (0-移动, 1-按下, 2-弹起), X Y 是指针位置, 最后两个参数似乎没有用到 (cpp 中直接获取的当前窗口尺寸) |

### 注意

1. 鼠标的 mouseState 非 Move 时才会处理 IsLeft。鼠标的延迟可能来源于前端 js 对用户输入的采样率
2. X Y 传入的是在 docker 中的光标位置，js 的光标位置传递给 go 时需要进行转换: `X = X * screenWidth / Width    Y = Y * screenHeight / Height`
3. 后台向 1xx09 端口被动建立的 tcp 连接发送 "K65,1|" 表示按下了 A 键

## 常见问题

1. 境内网络可能无法直接从 wine 源获取更新，我们使用指定 hosts 的方式来部分解决，若后续 ip 地址更换，需要手动替换 `docker build --add-host=dl.winehq.org:ipaddr` 

## ToDo

1. 更改过路径名，来自原项目的部分脚本内部需要认真修改
2. image 名称改为 dcwine，之前有使用 syncwine 的地方需要注意
3. container 进行了编号: "appvm${container_id}"，后台的控制需要针对性修改
4. 对屏幕尺寸进行控制
   1. Xvfb 控制屏幕尺寸
   2. ffmpeg 控制推流视频尺寸
   3. syncinput 控制应用窗口大小
5. syncinput 内部增加对 socket 的维护 (How to check if a socket pipe is broken?)
6. syncinput 对串行处理 string 的优化 (Use some proper serialization?)

## 潜在问题

1. 如果 syncinput 的连接丢失怎么办？如果无法重启连接，就相当于无任何按键响应
2. synvinput 针对 app 和 game 有不同的键盘输入模式，鼠标的输入都一样
   1. app 使用 virtual-key code
   2. game 使用 hardware scan code
3. 指定 docker depositary 时，需要将镜像名改为 ip:port/dcwine:latest，并且 depositary 最好提供 https，否则需要修改 docker 的配置文件
   1. 编辑 `/etc/docker/daemon.json`，添加 `"insecure-registries":["IP:Port"] `
   2. 重启 docker：`systemctl daemon-reload`，`systemctl restart docker`