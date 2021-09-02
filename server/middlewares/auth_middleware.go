package middlewares

import (
	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/repositories"
	"github.com/aliparlakci/country-roads/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func AuthMiddleware(userFinder repositories.UserFinder, sessions services.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		sessionId, err := c.Cookie("session")
		if err == http.ErrNoCookie {
			logger.Debug("request is anonymous")
			c.Next()
			return
		} else if err != nil {
			logger.Errorf("cannot extract session cookie from request headers: %v", err)
			c.Next()
			return
		}

		userId, err := sessions.FetchSession(c.Copy(), sessionId)
		if err != nil {
			logger.WithField("session_id", sessionId).Errorf("sessions.FetchSession raised an error when fetching session with session_id: %v", err.Error())
			c.Next()
			return
		}

		objId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			logger.WithField("user_id", userId).Errorf("cannot create objID from user_id: %v", err.Error())
			c.Next()
			return
		}

		user, err := userFinder.FindOne(c.Copy(), bson.M{"_id": objId})
		if err == mongo.ErrNoDocuments {
			logger.WithField("user_id", objId.Hex()).Debug("user with user_id does not exist")
			c.Next()
			return
		}
		if err != nil {
			logger.WithField("user_id", user.ID.Hex()).Errorf("UserFinder.FindOne() raised an error while finding user with user_id: %v", err.Error())
			c.Next()
			return
		}

		logger.WithField("email", user.Email).Debug("request belongs to user with email")
		c.Set("user", user)
		c.Next()
	}
}