# zookeeper 终端可视化客户端
![演示](screen1.gif)

## `-h`
传入 zookeeper 服务地址，不传默认为 `127.0.0.1:2181`

## 界面错乱解决办法

- windows
在 `cmd` 中执行命令：

```shell
chcp 65001
```

- mac or linux
终端执行命令：

```shell
export RUNEWIDTH_EASTASIAN=0
```
