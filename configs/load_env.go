package configs

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnvVar(path *string) {
	var err error
	if path == nil {
		path = new(string)
		*path = ".env"
	}
	err = godotenv.Load(*path)

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}
