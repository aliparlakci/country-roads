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
	db, close := common.InitilizeDb(db_uri)
	defer close()

	env := common.Env{
		Db: db,
	}

	router := gin.Default()
	v1 := router.Group("v1")

	controllers.RegisterRideController(v1, &env)

	router.Run(":8080")
}
