package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/thunthup/aimet-test/controllers"
)

func EventRoute(router *gin.Engine) {
	router.GET("/events", controllers.ListEvents)
	router.GET("/events/:id", controllers.GetEventById)
	router.POST("/events", controllers.CreateEvent)
	router.PUT("/events/:id", controllers.UpdateEvent)
	router.DELETE("/events/:id", controllers.DeleteEvent)

}
