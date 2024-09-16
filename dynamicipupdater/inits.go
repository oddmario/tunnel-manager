package dynamicipupdater

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initRoutes(engine *gin.Engine) {
	engine.GET("/get_pub_ip", getPublicIP)
	engine.POST("/update_ip", updateIPHandler)
}

func initErrors(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	engine.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "METHOD_NOT_ALLOWED", "message": "Method not allowed"})
	})
}
