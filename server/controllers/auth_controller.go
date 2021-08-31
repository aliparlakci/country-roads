package controllers

import (
	"github.com/aliparlakci/country-roads/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterAuthController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/auth/login", Login())
}
