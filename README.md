# Nft-Port
[![GitHub issues](https://img.shields.io/github/issues/FanhuaCloud/nft-port)](https://github.com/FanhuaCloud/nft-port/issues)
[![GitHub forks](https://img.shields.io/github/forks/FanhuaCloud/nft-port)](https://github.com/FanhuaCloud/nft-port/network)
[![GitHub stars](https://img.shields.io/github/stars/FanhuaCloud/nft-port)](https://github.com/FanhuaCloud/nft-port/stargazers)
![language](https://img.shields.io/badge/language-go-orange.svg)

基于nftables的端口中转工具
## 安装
您可以从GitHub上面下载预构建版本
```bash
wget https://github.com/FanhuaCloud/nft-port/releases/latest/download/nft-port.zip
unzip nft-port.zip
cd nft-port
chmod +x nft_port_amd64_linux
./nft_port_amd64_linux -a help
```
为了正常使用，您需要安装nftables，并且开启ip_forward
```bash
# 开启ip_forward
echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
```
## 构建
构建使用go build即可，全平台均可编译
## 使用
```bash
consoleLogger Init:{"level":"INFO","color":true,"LogLevel":0}
2020-03-10 20:03:03 [INFO] [nft-port/main/main.go:124] nft-port version 1.1
2020-03-10 20:03:03 [INFO] [nft-port/main/main.go:125] Aauthor: https://github.com/FanhuaCloud
  -a string
    	Actions that need to be performed, can use resolve, load, clear, list, nft. (default "help")
  -c string
    	config_path (default "./config.yaml")
  -d string
    	Domain names that need to be resolved (default "www.baidu.com")
```
### 测试HttpDNS
```bash
[root@iZbp15mr3banydZ ~]# ./nft_port_amd64_linux -a resolve -d www.baidu.com
consoleLogger Init:{"level":"DEBG","color":true,"LogLevel":0}
2020-03-10 16:50:27 [INFO] [nft-port/main/main.go:104] nft-port version 1.0
2020-03-10 16:50:27 [INFO] [nft-port/main/main.go:105] Aauthor: https://github.com/FanhuaCloud
2020-03-10 16:50:28 [INFO] [nft-port/main/main.go:56] 220.181.38.149;220.181.38.150
```
### 列出规则
```bash
./nft_port_amd64_linux -a list -c ./config.yaml
```
### 加载规则
```bash
[root@iZbp15mr3banydZ ~]# ./nft_port_amd64_linux -a load
consoleLogger Init:{"level":"DEBG","color":true,"LogLevel":0}
2020-03-10 16:51:47 [INFO] [nft-port/main/main.go:104] nft-port version 1.0
2020-03-10 16:51:47 [INFO] [nft-port/main/main.go:105] Aauthor: https://github.com/FanhuaCloud
2020-03-10 16:51:47 [INFO] [nft-port/main/main.go:61] Load config：./config.yaml
2020-03-10 16:51:47 [INFO] [nft-port/main/main.go:68] Gen the nft file to /tmp/ipv4-portforward.
2020-03-10 16:51:47 [INFO] [nft-port/main/yaml/yaml_util.go:54] Resolve domain： www.baidu.com
2020-03-10 16:51:47 [DEBG] [nft-port/main/yaml/yaml_util.go:68] #! /usr/sbin/nft -f

add table portforward
flush table ip portforward
add chain portforward prerouting { type nat hook prerouting priority -100; }
add chain portforward postrouting { type nat hook postrouting priority 100; }
add rule portforward postrouting mark 0x00000089 counter masquerade
add rule ip portforward prerouting tcp dport 1433 counter mark set 0x00000089 dnat to 1.1.1.1:443
add rule ip portforward prerouting tcp dport 1435 counter mark set 0x00000089 dnat to 0.0.0.0:443

2020-03-10 16:51:47 [INFO] [nft-port/main/main.go:33] Use nft -f to load rule.
2020-03-10 16:51:47 [INFO] [nft-port/main/main.go:47] Load rule successed.
```
### 清除规则
```bash
[root@iZbp15mr3banydZ ~]# ./nft_port_amd64_linux -a clear -c ./config.yaml 
consoleLogger Init:{"level":"DEBG","color":true,"LogLevel":0}
2020-03-10 16:53:00 [INFO] [nft-port/main/main.go:104] nft-port version 1.0
2020-03-10 16:53:00 [INFO] [nft-port/main/main.go:105] Aauthor: https://github.com/FanhuaCloud
2020-03-10 16:53:00 [INFO] [nft-port/main/main.go:74] Load config：./config.yaml
2020-03-10 16:53:00 [DEBG] [nft-port/main/yaml/yaml_util.go:77] #! /usr/sbin/nft -f

flush table ip portforward

2020-03-10 16:53:00 [INFO] [nft-port/main/main.go:33] Use nft -f to load rule.
2020-03-10 16:53:00 [INFO] [nft-port/main/main.go:47] Load rule successed.

```
### 查看原始nftables规则
```bash
[root@ecs-9JW nft-port]# ./nft_port_amd64_linux -a nft
consoleLogger Init:{"level":"INFO","color":true,"LogLevel":0}
2020-03-10 19:57:48 [INFO] [nft-port/main/main.go:124] nft-port version 1.1
2020-03-10 19:57:48 [INFO] [nft-port/main/main.go:125] Aauthor: https://github.com/FanhuaCloud
2020-03-10 19:57:48 [INFO] [nft-port/main/yaml/yaml_util.go:31] Load config：./config.yaml
table ip portforward {
	chain prerouting {
		type nat hook prerouting priority -100; policy accept;
	}

	chain postrouting {
		type nat hook postrouting priority 100; policy accept;
	}
}
```

## apiserver
使用-m启动apiserver，程序将以daemon方式接受rest请求，目前已经完成简单实现
```bash
[root@ecs-9JW ~]# ./nft_port_amd64_linux -m
consoleLogger Init:{"level":"INFO","color":true,"LogLevel":0}
2020-05-08 11:05:14 [INFO] [nft-port/main.go:26] nft-port version 1.1
2020-05-08 11:05:14 [INFO] [nft-port/main.go:27] Aauthor: https://github.com/FanhuaCloud
2020-05-08 11:05:14 [INFO] [nft-port/yaml/yaml_util.go:82] Load config：./config.yaml
2020-05-08 11:05:14 [INFO] [nft-port/main.go:43] Use the daemon mode
2020-05-08 11:05:14 [INFO] [nft-port/main.go:46] 8387
2020-05-08 11:05:14 [INFO] [nft-port/main.go:47] apikey
2020-05-08 11:05:14 [INFO] [nft-port/main.go:56] Start daemon api server
```
请求路由
```
r := gin.Default()
r.GET("/ping", ping)
v1 := r.Group("/v1")
v1.Use(checkApiKey(conf.DaemonConf.ApiKey))
{
    v1.DELETE("/port/:name", delPort)     // 删除port
    v1.PUT("/port/:name", addPort)        // 增加port
    v1.GET("/port/:name", getPort)        // 查看port信息
    v1.POST("/conf/reload", reloadConfig) // 重载配置
    v1.POST("/rules/reload", reloadRules) // 重载转发表
    v1.POST("/rules/clear", clearRules)   // 清空转发表
    v1.GET("/rules/list", listRules)      // 返回规则列表
}
```

## 配置文件
程序使用yaml格式配置文件，默认从./config.yaml读取，可使用 -c 指定配置文件

这是一个模板配置文件
```yaml
table-name: portforward

daemon:
  port: 8387
  key: "apikey"

port:
  - name: "test"
    type: ip # dns, or ip
    listen-port: 1433 # listen port
    server: server # server address
    port: 443 # server port
  - name: "test2"
    type: dns # dns, or ip
    listen-port: 1435 # listen port
    server: server # server address
    port: 443 # server port
```
## 感谢
- github.com/asmcos/requests
- github.com/wonderivan/logger
- gopkg.in/yaml.v3
- github.com/gin-gonic/gin