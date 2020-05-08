package yaml

import (
	"bufio"
	"fmt"
	httpDns "github.com/FanhuaCloud/nft-port/dns"
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	TableName string `yaml:"table-name"`
	Port      []Port `yaml:"port"`
}

type Port struct {
	//- name: "test"
	//type: dns # dns, or ip
	//listen-port: 1433 # listen port
	//server: server # server address
	//port: 443 # server port
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	ListenPort int    `yaml:"listen-port"`
	Server     string `yaml:"server"`
	ServerPort int    `yaml:"port"`
}

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

func ReadYaml(configPath string) (*Config, error) {
	logger.Info("Load config：" + configPath)
	conf := new(Config)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (cfg *Config) GenRule() string {
	table := cfg.TableName
	var nftCMD = new(strings.Builder)
	nftCMD.WriteString("#! /usr/sbin/nft -f\n\n")
	nftCMD.WriteString(fmt.Sprintf("add table %s\n", table))
	nftCMD.WriteString(fmt.Sprintf("flush table ip %s\n", table))
	nftCMD.WriteString(fmt.Sprintf("add chain %s prerouting { type nat hook prerouting priority -100; }\n", table))
	nftCMD.WriteString(fmt.Sprintf("add chain %s postrouting { type nat hook postrouting priority 100; }\n", table))
	nftCMD.WriteString(fmt.Sprintf("add rule %s postrouting mark 0x00000089 counter masquerade\n", table))
	for _, e := range cfg.Port {
		if e.Type == "dns" {
			logger.Info("Resolve domain：", e.Server)
			ip, err := httpDns.Resolve(e.Server)
			if err != nil {
				logger.Error(err)
				continue
			}
			if ip == "" {
				logger.Error("Not a vaild domain：", e.Server)
				continue
			}
			e.Server = strings.Split(ip, ";")[0]
		}
		nftCMD.WriteString(fmt.Sprintf("add rule ip %s prerouting tcp dport %d counter mark set 0x00000089 dnat to %s:%d\n", table, e.ListenPort, e.Server, e.ServerPort))
		nftCMD.WriteString(fmt.Sprintf("add rule ip %s prerouting udp dport %d counter mark set 0x00000089 dnat to %s:%d\n", table, e.ListenPort, e.Server, e.ServerPort))
	}
	logger.Debug(nftCMD.String())
	return nftCMD.String()
}

func (cfg *Config) GenClearRule() string {
	table := cfg.TableName
	var nftCMD = new(strings.Builder)
	nftCMD.WriteString("#! /usr/sbin/nft -f\n\n")
	nftCMD.WriteString(fmt.Sprintf("flush table ip %s\n", table))
	logger.Debug(nftCMD.String())
	return nftCMD.String()
}

func (cfg *Config) ListNftRules() {
	//nft list table ip portforward
	cmd := exec.Command("nft", "list", "table", cfg.TableName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run 和 Start只能用一个
	err := cmd.Run()
	if err != nil {
		logger.Error(err)
	}
}

func (cfg *Config) ListRules() {
	logger.Info("Name   ServerPort   Server   ListenPort   Type")
	for _, e := range cfg.Port {
		logger.Info(e.Name, "   ", e.ServerPort, "   ", e.Server, "   ", e.ListenPort, "   ", e.Type)
	}
}

func (cfg *Config) LoadRules() {
	//生成规则文件
	logger.Info("Gen the nft file to /tmp/ipv4-portforward.")
	runRules(cfg.GenRule())
}

func (cfg *Config) ClearRules() {
	runRules(cfg.GenClearRule())
}
