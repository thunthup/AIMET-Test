package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/thunthup/aimet-test/controllers"
)

func EventRoute(router *gin.Engine) {
	router.GET("/api/events", controllers.ListEvents)
	router.GET("/api/events/:id", controllers.GetEventById)
	router.POST("/api/events", controllers.CreateEvent)
	router.PUT("/api/events/:id", controllers.UpdateEvent)
	router.DELETE("/api/events/:id", controllers.DeleteEvent)

}
