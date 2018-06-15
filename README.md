# 码云 webhook
## config.json 配置：

- ListenPort：监听端口
- RepoName：url path，程序根据RepoName匹配WorkerPath
- WorkPath：git仓库路径
- Secret：webhook密码

## Usage
1 .下载
```bash
go get github.com/lishuhao/webhook
```
2 .编译（Linux运行环境）
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 【-o 指定程序名】 main.go
```
3 . cp conf.simple.json conf.json

4 . 后台运行
```bash
command >out.file 2>&1 &
```
- command>out.file是将command的输出重定向到out.file文件，即输出内容不打印到屏幕上，而是输出到out.file文件中。
- 2>&1 是将标准出错重定向到标准输出，这里的标准输出已经重定向到了out.file文件，即将标准出错也输出到out.file文件中。最后一个&， 是让该命令在后台执行。

## linux 可执行程序
webhook_linux