package controllers

import (
	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateUser(users repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: Implement
		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterUserController(router *gin.RouterGroup, env *common.Env) {
	router.PUT("/users/:id", UpdateUser(env.Repositories.UserRepository))
}
