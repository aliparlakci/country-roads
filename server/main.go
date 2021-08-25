package main

import (
	"os"

	"example.com/country-roads/common"
	"example.com/country-roads/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var env *common.Env
	{
		db_uri := os.Getenv("DB_CONNECTION")
		db, close := common.InitilizeDb(db_uri)
		defer close()

		rdb_uri := os.Getenv("REDIS_ADDR")
		rdb := common.InitilizeRedis(rdb_uri, "", 0)

		env = &common.Env{
			Db:  db,
			Rdb: rdb,
		}
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("futuredate", futureDate)
		v.RegisterValidation("validridetype", validRideType)
		v.RegisterValidation("validdirection", validDirection)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	api := router.Group("api")
	controllers.RegisterRideController(api, env)

	router.Run(":8080")
}
