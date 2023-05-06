package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thunthup/aimet-test/configs"
	"github.com/thunthup/aimet-test/routers"
)

func init() {
	configs.LoadEnvVar()
	configs.ConnectPostgresDB()
}

func main() {
	router := gin.New()
	routers.HealthCheckRoute(router)
	routers.EventRoute(router)
	fmt.Println("server is running on", os.Getenv("PORT"))
	router.Run()
}
