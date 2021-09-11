package controllers

import (
	"fmt"
	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/middlewares"
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func GetContact(finder repositories.ContactFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c)

		u, _ := c.Get("user")
		user := u.(models.User)

		id := c.Param("id")
		userId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			logger.WithField("user_id", c.Param("id")).Debug("user_id is not a valid object ID")
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		contact, err := finder.FindOne(c.Copy(), bson.M{"owner": userId})

		if err == mongo.ErrNoDocuments {
		} else if err != nil {
			logger.WithField("user_id", userId.Hex()).Errorf("ContactFinder.FindOne raised an error while fetching the document _id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, models.ContactResponse{
			Name:     user.DisplayName,
			Email:    user.Email,
			Phone:    contact.Phone,
			Whatsapp: contact.Whatsapp,
		})
	}
}

func PutContact(contacts repositories.ContactRepository, users repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c)

		u, exists := c.Get("user")
		if !exists {
			logger.Errorf("user does not exists in the context")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		user := u.(models.User)

		var contactDto models.UpdateContactForm
		if err := contactDto.Bind(c); err != nil {
			logger.Debugf("cannot bind to models.UpdateContactForm: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Contact format was incorrect: %v", err)})
			return
		}

		if _, err := contacts.FindOne(c.Copy(), bson.M{"owner": user.ID}); err == mongo.ErrNoDocuments { // Contact document does exist, create a new one
			logger.WithField("email", user.Email).Debug("user with email has no contact information registered")

			if _, err := contacts.InsertOne(c.Copy(), models.Contact{
				Owner:    user.ID,
				Phone:    contactDto.Phone,
				Whatsapp: contactDto.Whatsapp,
			}); err != nil {
				logger.Errorf("ContactRepository.InsertOne raised an error: %v", err)
				return
			}

		} else if err != nil {
			logger.Errorf("ContactRepository.FindOne raised an error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		} else { // Contact document exists, update the existing one

			if err := contacts.UpdateOne(c.Copy(), bson.M{"owner": user.ID}, bson.D{
				{"$set", bson.D{
					{"phone", contactDto.Phone},
					{"whatsapp", contactDto.Whatsapp},
				}},
			}); err != nil {
				logger.WithField("owner", user.ID.Hex()).Errorf("ContactRepository.UpdateOne raised an error while updating document with owner: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

		}

		if err := users.UpdateOne(c.Copy(), bson.M{"_id": user.ID}, bson.D{
			{"$set", bson.D{
				{"displayName", contactDto.DisplayName},
			}},
		}); err != nil {
			logger.
				WithField("user_id", user.ID.Hex()).
				WithField("displayName", contactDto.DisplayName).
				Errorf("UserRepository.UpdateOne raised an error while updating the displayName of user: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterContactController(router *gin.RouterGroup, env *common.Env) {
	router.GET("/contact/:id", middlewares.Protected(GetContact(env.Repositories.ContactRepository)))
	router.PUT("/contact", middlewares.Protected(PutContact(
		env.Repositories.ContactRepository,
		env.Repositories.UserRepository,
	)))
}
