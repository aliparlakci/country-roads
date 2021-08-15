package main

import (
	"os"

	"example.com/country-roads/common"
	"example.com/country-roads/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db_uri := os.Getenv("DB_CONNECTION")

	var env common.Env
	{ // Prepare dependencies
		client, close := common.InitilizeDb(db_uri)
		env.Db = client
		defer close()
	}

	router := gin.Default()
	v1 := router.Group("v1")

	controllers.RegisterRideController(v1, &env)

	router.Run(":8080")
}
