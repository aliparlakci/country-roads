package controllers

import (
	"example.com/country-roads/interfaces"
	"fmt"
	"net/http"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getAllLocations(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		results, err := models.GetLocations(ctx, env.Db)
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

		var validator interfaces.Validator = locationDto
		if isValid, err := validator.Validate(); !isValid || err != nil {
			ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Location format was invalid: %v", err))
		}

		if locationDto.ParentID != "" {
			objID, err := primitive.ObjectIDFromHex(locationDto.ParentID)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Location format was invalid: %v", err))
			}

			if _, err := models.GetSingleLocation(ctx, env.Db, objID); err != nil {
				ctx.JSON(http.StatusBadRequest, "Location format was invalid")
			}
		}

		id, err := models.RegisterLocation(ctx, env.Db, locationDto)
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
