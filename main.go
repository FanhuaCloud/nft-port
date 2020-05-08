package main

import (
	"flag"
	"fmt"
	"github.com/FanhuaCloud/nft-port/api"
	httpDns "github.com/FanhuaCloud/nft-port/dns"
	yamlUtil "github.com/FanhuaCloud/nft-port/yaml"
	"github.com/wonderivan/logger"
)

func resolveDomain(domain *string) {
	ip, err := httpDns.Resolve(*domain)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(ip)
}

func main() {
	//初始化日志
	err := logger.SetLogger(`{"Console": {"level": "INFO","color": true}}`)
	if err != nil {
		logger.Error(err)
	}
	//输出版本信息
	logger.Info("nft-port version 1.2")
	logger.Info("Aauthor: https://github.com/FanhuaCloud")

	//设置flag
	action := flag.String("a", "help", "Actions that need to be performed, can use resolve, load, clear, list, nft.")
	domain := flag.String("d", "www.baidu.com", "Domain names that need to be resolved")
	configPath := flag.String("c", "./config.yaml", "config_path")
	isDaemon := flag.Bool("m", false, "Use daemon mode")
	flag.Parse()

	conf, err := yamlUtil.ReadYaml(*configPath)
	if err != nil {
		logger.Error("Read config failed.")
		return
	}

	if *isDaemon {
		logger.Info("Use the daemon mode")
		daemonListenIP := conf.DaemonConf.ListenIP
		daemonPort := conf.DaemonConf.Port
		if daemonPort == 0 {
			daemonPort = 8766
		}
		// 判断端口
		if daemonPort <= 0 || daemonPort > 65535 {
			logger.Error("Daemon port is error.")
			return
		}
		logger.Info("Start daemon api server")
		err := api.RunApiServer(fmt.Sprintf("%s:%d", daemonListenIP, daemonPort), conf)
		if err != nil {
			logger.Error(err)
		}
		return
	}

	//解析action
	switch *action {
	case "resolve":
		//http_dns解析测试
		resolveDomain(domain)
		break
	case "load":
		//加载规则
		conf.LoadRules()
		break
	case "clear":
		//清除规则
		conf.ClearRules()
		break
	case "list":
		// 列出所有规则
		conf.ListRules()
	case "nft":
		// 查看nft规则
		conf.ListNftRules()
	default:
		flag.PrintDefaults()
		break
	}
}
