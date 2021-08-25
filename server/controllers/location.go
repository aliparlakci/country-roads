package controllers

import (
	"net/http"

	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"github.com/gin-gonic/gin"
)

func getAllLocations(env *common.Env) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		results, err := models.GetLocations(ctx, env.Db.Database("country-roads"))
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, results)
	}
}

func RegisterLocationController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/locations", getAllLocations(env))
}
