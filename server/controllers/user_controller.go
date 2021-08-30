package controllers

import (
	"example.com/country-roads/common"
	"example.com/country-roads/models"
	"example.com/country-roads/validators"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SignUp(finder models.UserFinder, validatorFactory validators.IValidatorFactory) gin.HandlerFunc {
	validator, err := validatorFactory.GetValidator("users")
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		validator.Validate(c)
		c.JSON(http.StatusOK, gin.H{})
	}
}

func UpdateDisplayName(findUpdater models.UserFindUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func UpdatePhone(findUpdater models.UserFindUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterUserController(router *gin.RouterGroup, env *common.Env) {
	router.POST("/users/", SignUp(
			env.Repositories.UserRepository,
			env.ValidatorFactory,
		))
	router.PUT("/users/name", UpdateDisplayName(env.Repositories.UserRepository))
	router.PUT("/users/phone", UpdatePhone(env.Repositories.UserRepository))
}