package main

import (
	"github.com/aliparlakci/country-roads/middlewares"
	"github.com/aliparlakci/country-roads/services"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/controllers"
	"github.com/aliparlakci/country-roads/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{PadLevelText: true})
	//logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})

	env := &common.Env{}
	{
		dbUri := os.Getenv("MDB_URI")
		dbName := "country-roads"
		db, close := common.InitializeDb(dbUri, dbName)
		defer close()

		rdbUri := os.Getenv("RDB_URI")
		env.Repositories = &common.RepositoryContainer{
			RideRepository:     &models.RideCollection{Collection: db.Collection("rides")},
			LocationRepository: &models.LocationCollection{Collection: db.Collection("locations")},
			UserRepository:     &models.UserCollection{Collection: db.Collection("users")},
		}
		env.Services = &common.ServiceContainer{
			SessionService: &services.SessionStore{Store: common.InitializeRedis(rdbUri, "", 0)},
			OTPService:     &services.OTPStore{Store: common.InitializeRedis(rdbUri, "", 1)},
		}
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("request_id", uuid.New().String())
		c.Next()
	})
	router.Use(middlewares.Logger())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE"},
	}))

	router.Use(middlewares.AuthMiddleware(
		env.Repositories.UserRepository,
		env.Services.SessionService,
	))
	api := router.Group("api")
	{
		controllers.RegisterRideController(api, env)
		controllers.RegisterLocationController(api, env)
		controllers.RegisterUserController(api, env)
		controllers.RegisterAuthController(api, env)
	}

	router.Run(":8080")
}
