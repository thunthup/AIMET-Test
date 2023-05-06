package routers

import "github.com/gin-gonic/gin"

func HealthCheckRoute(router *gin.Engine) {
	router.GET("/healthcheck", func(c *gin.Context) {
		c.String(200, "ok")
	})

}
