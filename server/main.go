package main

import (
	"example.com/country-roads/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"

	"example.com/country-roads/models"
	"example.com/country-roads/validators"

	"example.com/country-roads/common"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	env := &common.Env{}
	{
		dbUri := os.Getenv("MDB_URI")
		dbName := "country-roads"
		db, close := common.InitilizeDb(dbUri, dbName)
		defer close()

		rdbUri := os.Getenv("RDB_URI")
		rdb := common.InitilizeRedis(rdbUri, "", 0)

		env.Rdb = rdb
		env.Collections = &common.CollectionContainer{
			RideCollection:     &models.RideCollection{Collection: db.Collection("rides")},
			LocationCollection: &models.LocationCollection{Collection: db.Collection("locations")},
		}
		env.Validators = &validators.ValidatorFactory{LocationFinder: env.Collections.LocationCollection}
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE"},
	}))

	api := router.Group("api")
	controllers.RegisterRideController(api, env)
	controllers.RegisterLocationController(api, env)

	router.Run(":8080")
}
