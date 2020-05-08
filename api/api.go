package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type apiHeader struct {
	ApiKey string `header:"apikey"`
}

func checkApiKey(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := apiHeader{}
		err := c.ShouldBindHeader(&h)
		if err != nil {
			_ = c.AbortWithError(500, err)
			return
		}
		if h.ApiKey != apiKey {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func delPort(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func addPort(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func getPort(c *gin.Context) {
	c.String(http.StatusOK, "pong")
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
