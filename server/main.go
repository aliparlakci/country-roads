package main

import (
	"example.com/country-roads/interfaces"
	"example.com/country-roads/models"
	"example.com/country-roads/validators"
	"os"

	"example.com/country-roads/common"
	"example.com/country-roads/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var env *common.Env
	{
		dbUri := os.Getenv("MDB_URI")
		dbName := "country-roads"
		db, close := common.InitilizeDb(dbUri, dbName)
		defer close()

		rdbUri := os.Getenv("RDB_URI")
		rdb := common.InitilizeRedis(rdbUri, "", 0)

		env = &common.Env{
			Collections: common.CollectionContainer{
				RideCollection:     &models.RideCollection{Collection: db.Collection("rides")},
				LocationCollection: &models.LocationCollection{Collection: db.Collection("locations")},
			},
			Validators: common.ValidatorContainer{
				RideValidator:     func() interfaces.Validator {
					return &validators.RideValidator{LocationFinder: env.Collections.LocationCollection}
				},
				LocationValidator: func() interfaces.Validator {
					return &validators.LocationValidator{}
				},
			},
			Rdb: rdb,
		}
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
