package middlewares

import (
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(finder models.UserFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionContext, exists := c.Get("session")
		if !exists {
			c.Next()
			return
		}

		session, success := sessionContext.(services.Session)
		if !success {
			c.Next()
			return
		}

		userId := session.UserId
		if userId == "" {
			c.Next()
			return
		}

		objId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			//TODO: Log error here
			c.Next()
			return
		}

		user, err := finder.FindOne(c.Copy(), bson.M{"_id": objId})
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}