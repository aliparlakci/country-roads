package controllers

import (
	"fmt"
	"net/http"

	"example.com/country-roads/validators"
	"go.mongodb.org/mongo-driver/bson"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
)

func GetAllLocations(finder models.LocationFinder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		results, err := finder.FindMany(ctx, bson.D{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"results": results})
	}
}

func PostLocation(inserter models.LocationInserter, validators validators.IValidatorFactory) gin.HandlerFunc {
	validator, err := validators.GetValidator("locations")
	if err != nil {
		panic(err)
	}
	return func(ctx *gin.Context) {
		var locationDto models.NewLocationForm

		if err := ctx.Bind(&locationDto); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Location format was incorrect: %v", err)})
			return
		}

		validator.SetDto(locationDto)
		if isValid, err := validator.Validate(ctx); !isValid || err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Location format was invalid"})
			return
		}

		var schema models.LocationSchema
		if locationDto.ParentKey != "" {
			schema = models.LocationSchema{Key: locationDto.Key, Display: locationDto.Display, ParentKey: locationDto.ParentKey}
		} else {
			schema = models.LocationSchema{Key: locationDto.Key, Display: locationDto.Display}
		}

		id, err := inserter.InsertOne(ctx, schema)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Location couldn't get created: %v", err)})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func RegisterLocationController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/locations", GetAllLocations(env.Repositories.LocationRepository))
	router.POST("/locations", PostLocation(
		env.Repositories.LocationRepository,
		env.ValidatorFactory,
	))
}
