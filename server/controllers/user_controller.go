package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/validators"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func PostUser(findInserter models.UserFindInserter, validatorFactory validators.IValidatorFactory) gin.HandlerFunc {
	validator, err := validatorFactory.GetValidator("users")
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		var userDto models.NewUserForm

		if err := c.Bind(&userDto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user format was incorrect"})
			return
		}

		if err := validator.SetDto(&userDto); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if result, err := validator.Validate(c.Copy()); !result || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is not valid"})
			return
		}
		if _, err := findInserter.FindOne(c.Copy(), bson.M{"email": userDto.Email}); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}
		id, err := findInserter.InsertOne(c.Copy(), models.UserSchema{
			DisplayName: userDto.DisplayName,
			Email:       userDto.Email,
			Phone:       userDto.Phone,
			Verified:    false,
			SignedUpAt:  time.Now(),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("new user could not get created: %v", err)})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func UpdateDisplayName(findUpdater models.UserFindUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: Implement
		c.JSON(http.StatusOK, gin.H{})
	}
}

func UpdatePhone(findUpdater models.UserFindUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: Implement
		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterUserController(router *gin.RouterGroup, env *common.Env) {
	router.POST("/users", PostUser(
		env.Repositories.UserRepository,
		env.ValidatorFactory,
	))
	// router.PUT("/users/name", UpdateDisplayName(env.Repositories.UserRepository))
	// router.PUT("/users/phone", UpdatePhone(env.Repositories.UserRepository))
}
