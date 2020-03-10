package main

import (
	"bufio"
	"flag"
	"github.com/wonderivan/logger"
	httpDns "nft-port/main/dns"
	yamlUtil "nft-port/main/yaml"
	"os"
	"os/exec"
)

func runRules(rule string) {
	// 写入
	fi, err := os.Create("/tmp/ipv4-portforward")
	if err != nil {
		logger.Error(err)
		return
	}
	defer fi.Close()
	w := bufio.NewWriter(fi)
	_, err = w.WriteString(rule)
	if err != nil {
		logger.Error(err)
		return
	}
	//判断是否写入成功
	if err = w.Flush(); err != nil {
		logger.Error(err)
		return
	}
	//开始加载规则
	logger.Info("Use nft -f to load rule.")
	_, err = exec.LookPath("nft")
	//查找nft，不存在报错
	if err != nil {
		logger.Error(err)
		return
	}
	cmd := exec.Command("nft", "-f", "/tmp/ipv4-portforward")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run 和 Start只能用一个
	err = cmd.Run()
	if err != nil {
		logger.Error(err)
	}
	if !cmd.ProcessState.Success() {
		logger.Info("Load rule failed, please check the stderr.")
	} else {
		logger.Info("Load rule successed.")
	}
}

func resolveDomain(domain *string) {
	ip, err := httpDns.Resolve(*domain)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(ip)
}

func loadRules(configPath *string) {
	//加载yaml
	conf, err := yamlUtil.ReadYaml(*configPath)
	if err != nil {
		logger.Error(err)
		return
	}
	//生成规则文件
	logger.Info("Gen the nft file to /tmp/ipv4-portforward.")
	runRules(conf.GenRule())
}

func clearRules(configPath *string) {
	//加载yaml
	conf, err := yamlUtil.ReadYaml(*configPath)
	if err != nil {
		logger.Error(err)
		return
	}
	runRules(conf.ClearRule())
}

func listRules(configPath *string) {
	//加载yaml
	conf, err := yamlUtil.ReadYaml(*configPath)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("Name   ServerPort   Server   ListenPort   Type")
	for _, e := range conf.Port {
		logger.Info(e.Name, "   ", e.ServerPort, "   ", e.Server, "   ", e.ListenPort, "   ", e.Type)
	}
}

func main() {
	//初始化日志
	err := logger.SetLogger(`{"Console": {"level": "INFO","color": true}}`)
	if err != nil {
		logger.Error(err)
	}
	//输出版本信息
	logger.Info("nft-port version 1.1")
	logger.Info("Aauthor: https://github.com/FanhuaCloud")

	//设置flag
	action := flag.String("a", "help", "Actions that need to be performed, can use resolve, load, clear, list.")
	domain := flag.String("d", "www.baidu.com", "Domain names that need to be resolved")
	configPath := flag.String("c", "./config.yaml", "config_path")
	flag.Parse()

	//解析action
	switch *action {
	case "resolve":
		//http_dns解析测试
		resolveDomain(domain)
		break
	case "load":
		//加载规则
		loadRules(configPath)
		break
	case "clear":
		//清除规则
		clearRules(configPath)
		break
	case "list":
		// 列出所有规则
		listRules(configPath)
	default:
		flag.PrintDefaults()
		break
	}
}
