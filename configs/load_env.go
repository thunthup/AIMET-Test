package configs

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnvVar() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}
