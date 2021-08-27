package controllers

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"example.com/country-roads/schemas"
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

		if len(results) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"results": results})
	}
}

func PostLocation(inserter models.LocationInserter, getValidator func() validators.Validator) gin.HandlerFunc {
	validator := getValidator()
	return func(ctx *gin.Context) {
		var locationDto models.LocationDTO

		if err := ctx.Bind(&locationDto); err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Location format was incorrect: %v", err))
			return
		}

		validator.SetDto(locationDto)
		if isValid, err := validator.Validate(ctx); !isValid || err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Location format was invalid: %v", err))
		}

		var schema schemas.LocationSchema
		if locationDto.ParentID != "" {
			parentId, _ := primitive.ObjectIDFromHex(locationDto.ParentID)
			schema = schemas.LocationSchema{Display: locationDto.Display, ParentID: parentId}
		} else {
			schema = schemas.LocationSchema{Display: locationDto.Display}
		}

		id, err := inserter.InsertOne(ctx, schema)
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("Location couldn't get created: %v", err))
			return
		}

		ctx.JSON(http.StatusCreated, id)
	}
}

func RegisterLocationController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/locations", GetAllLocations(env.Collections.LocationCollection))
	router.POST("/locations", PostLocation(
		env.Collections.LocationCollection,
		env.Validators.LocationValidator,
	))
}
