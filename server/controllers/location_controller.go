package controllers

import (
	"fmt"
	"net/http"

	"example.com/country-roads/schemas"
	"example.com/country-roads/validators"
	"go.mongodb.org/mongo-driver/bson"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getAllLocations(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var finder models.LocationFinder = env.Collections.LocationCollection
		results, err := finder.FindMany(ctx, bson.D{})
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, results)
	}
}

func postLocation(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var locationDto models.LocationDTO

		if err := ctx.Bind(&locationDto); err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Location format was incorrect: %v", err))
			return
		}

		var validator validators.Validator = env.Validators.LocationValidator()
		validator.SetDto(locationDto)
		if isValid, err := validator.Validate(ctx); !isValid || err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Location format was invalid: %v", err))
		}

		var schema schemas.LocationSchema
		if locationDto.ParentID != "" {
			parentId, err := primitive.ObjectIDFromHex(locationDto.ParentID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Location format was invalid: %v", err))
			}

			var finder models.LocationFinder = env.Collections.LocationCollection
			if _, err := finder.FindOne(ctx, bson.M{"_id": parentId}); err != nil {
				ctx.JSON(http.StatusBadRequest, "Location format was invalid")
			}

			schema = schemas.LocationSchema{Display: locationDto.Display, ParentID: parentId}
		} else {
			schema = schemas.LocationSchema{Display: locationDto.Display}
		}

		var inserter models.LocationInserter = env.Collections.LocationCollection
		id, err := inserter.InsertOne(ctx, schema)
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Location couldn't get created: %v", err))
			return
		}

		ctx.JSON(http.StatusCreated, id)
	}
}

func RegisterLocationController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/locations", getAllLocations(env))
	router.POST("/locations", postLocation(env))
}
