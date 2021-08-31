package main

import (
	"github.com/aliparlakci/country-roads/middlewares"
	"github.com/aliparlakci/country-roads/services"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/controllers"
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/validators"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	env := &common.Env{}
	{
		dbUri := os.Getenv("MDB_URI")
		dbName := "country-roads"
		db, close := common.InitializeDb(dbUri, dbName)
		defer close()

		rdbUri := os.Getenv("RDB_URI")
		rdb := common.InitializeRedis(rdbUri, "", 0)

		env.Repositories = &common.RepositoryContainer{
			RideRepository:     &models.RideCollection{Collection: db.Collection("rides")},
			LocationRepository: &models.LocationCollection{Collection: db.Collection("locations")},
			UserRepository:     &models.UserCollection{Collection: db.Collection("users")},
		}
		env.Services = &common.ServiceContainer{
			SessionService: &services.SessionsClient{Client: rdb},
		}
		env.ValidatorFactory = &validators.ValidatorFactory{LocationFinder: env.Repositories.LocationRepository}
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE"},
	}))
	router.Use(middlewares.SessionMiddleware(env.Services.SessionService))
	router.Use(middlewares.AuthMiddleware(env.Repositories.UserRepository))
	api := router.Group("api")
	{
		controllers.RegisterRideController(api, env)
		controllers.RegisterLocationController(api, env)
		controllers.RegisterUserController(api, env)
		controllers.RegisterAuthController(api, env)
	}

	router.Run(":8080")
}
