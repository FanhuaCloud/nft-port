package api

import (
	"bufio"
	yamlUtil "github.com/FanhuaCloud/nft-port/yaml"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

type apiHeader struct {
	ApiKey string `header:"apikey"`
}

func checkApiKey(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := apiHeader{}
		err := c.ShouldBindHeader(&h)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if h.ApiKey != apiKey {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

func portExist(name string) (yamlBuild []yamlUtil.Port) {
	for _, e := range conf.Port {
		if e.Name == name {
			yamlBuild = append(yamlBuild, e)
		}
	}
	return yamlBuild
}

func writeFile(value string) error {
	fi, err := os.Create(confPath)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer fi.Close()
	w := bufio.NewWriter(fi)
	_, err = w.WriteString(value)
	if err != nil {
		logger.Error(err)
		return err
	}
	//判断是否写入成功
	if err = w.Flush(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func delPort(c *gin.Context) {
	name := c.Param("name")
	ports := portExist(name)
	if ports == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// 生成新的yaml
	var yamlBuild []yamlUtil.Port
	for _, e := range conf.Port {
		if e.Name != name {
			yamlBuild = append(yamlBuild, e)
		}
	}
	conf.Port = yamlBuild
	// 处理
	d, err := yaml.Marshal(&conf)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = writeFile(string(d))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	conf.LoadRules()
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func addPort(c *gin.Context) {
	port := yamlUtil.Port{}
	err := c.ShouldBindJSON(&port)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	conf.Port = append(conf.Port, port)
	port.InstallRules(conf.TableName)
	d, err := yaml.Marshal(&conf)
	if err != nil {
		logger.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	err = writeFile(string(d))
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func getPort(c *gin.Context) {
	ports := portExist(c.Param("name"))
	if ports == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   ports,
	})
}

func reloadRules(c *gin.Context) {
	conf.LoadRules()
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func clearRules(c *gin.Context) {
	conf.ClearRules()
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func listRules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   conf.Port,
	})
}

func reloadConfig(c *gin.Context) {
	confRead, err := yamlUtil.ReadYaml(confPath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	conf = confRead
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}
