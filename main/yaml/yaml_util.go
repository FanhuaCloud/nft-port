package yaml

import (
	"fmt"
	httpDns "github.com/FanhuaCloud/nft-port/main/dns"
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func (port *Config) GenRule() string {
	table := port.TableName
	var nftCMD = new(strings.Builder)
	nftCMD.WriteString("#! /usr/sbin/nft -f\n\n")
	nftCMD.WriteString(fmt.Sprintf("add table %s\n", table))
	nftCMD.WriteString(fmt.Sprintf("flush table ip %s\n", table))
	nftCMD.WriteString(fmt.Sprintf("add chain %s prerouting { type nat hook prerouting priority -100; }\n", table))
	nftCMD.WriteString(fmt.Sprintf("add chain %s postrouting { type nat hook postrouting priority 100; }\n", table))
	nftCMD.WriteString(fmt.Sprintf("add rule %s postrouting mark 0x00000089 counter masquerade\n", table))
	for _, e := range port.Port {
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

func (port *Config) ClearRule() string {
	table := port.TableName
	var nftCMD = new(strings.Builder)
	nftCMD.WriteString("#! /usr/sbin/nft -f\n\n")
	nftCMD.WriteString(fmt.Sprintf("flush table ip %s\n", table))
	logger.Debug(nftCMD.String())
	return nftCMD.String()
}
