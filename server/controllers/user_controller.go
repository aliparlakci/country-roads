package controllers

import (
	"fmt"
	"github.com/aliparlakci/country-roads/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func PostUser(findInserter repositories.UserFindInserter) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var userDto models.NewUserForm
		if err := userDto.Bind(c); err != nil {
			logger.WithField("body", c.Request.MultipartForm).Debug("cannot bind request body to models.NewUserForm")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user format was incorrect"})
			return
		}

		if _, err := findInserter.FindOne(c.Copy(), bson.M{"email": userDto.Email}); err == nil {
			logger.WithField("email", userDto.Email).Debug("user with email already exists")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
			return
		}
		id, err := findInserter.InsertOne(c.Copy(), models.UserSchema{
			DisplayName:      userDto.DisplayName,
			Email:            userDto.Email,
			Phone:            userDto.Phone,
			Verified:         false,
			SignedUpAt:       time.Now(),
			ContactWithPhone: userDto.ContactWithPhone,
		})
		if err != nil {
			logger.Errorf("models.UserFindInserter.InsertOne() raised an error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("new user could not get created: %v", err)})
			return
		}
		logger.WithField("id", id).Info("new user with id is created")
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func GetContact(users repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c)

		var authUser models.User
		if u, exists := c.Get("user"); !exists {
			logger.Debug("request is not authorized")
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			authUser = u.(models.User)
		}

		id := c.Param("id")
		userId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			logger.WithField("user_id", c.Param("id")).Debug("user_id is not a valid object ID")
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		user, err := users.FindOne(c.Copy(), bson.M{"_id": userId})
		if err == mongo.ErrNoDocuments {
			logger.WithField("user_id", userId.Hex()).Debug("user with user_id does not exist")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
			return
		} else if err != nil {
			logger.WithField("user_id", userId.Hex()).Errorf("UserRepository.FindOne raised an error while fetching the user with id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		response := models.ContactResponse{Email: user.Email}
		if !user.ContactWithPhone {
			if authUser.ID != user.ID {
				logger.WithField("user_id", user.ID.Hex()).Debug("user is not authorized to view the phone number")
				c.JSON(http.StatusOK, response)
				return
			}
		}

		response.Phone = user.Phone
		c.JSON(http.StatusOK, response)
	}
}

func UpdateUser(findUpdater repositories.UserFindUpdater) gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: Implement
		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterUserController(router *gin.RouterGroup, env *common.Env) {
	router.POST("/users", PostUser(env.Repositories.UserRepository))
	router.GET("/users/contact/:id", GetContact(env.Repositories.UserRepository))
	router.PUT("/users/:id", UpdateUser(env.Repositories.UserRepository))
	// router.PUT("/users/phone", UpdatePhone(env.Repositories.UserRepository))
}
