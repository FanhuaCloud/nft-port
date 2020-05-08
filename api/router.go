package api

import (
	"github.com/FanhuaCloud/nft-port/yaml"
	"github.com/gin-gonic/gin"
)

var conf *yaml.Config

func RunApiServer(listen string, cfg *yaml.Config) error {
	conf = cfg
	r := gin.Default()
	r.GET("/ping", ping)
	v1 := r.Group("/v1")
	v1.Use(checkApiKey(conf.DaemonConf.ApiKey))
	{
		v1.DELETE("/port/:name", delPort)     // 删除port
		v1.PUT("/port/:name", addPort)        // 增加port
		v1.GET("/port/:name", getPort)        // 查看port信息
		v1.POST("/rules/reload", reloadRules) // 重载转发表
		v1.POST("/rules/clear", clearRules)   // 清空转发表
		v1.GET("/rules/list", listRules)      // 返回规则列表
	}
	return r.Run(listen)
}
