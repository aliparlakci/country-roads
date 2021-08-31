package controllers

import (
	"fmt"
	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Login(finder models.UserFinder, sessions services.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds models.LoginRequestForm
		if err := c.Bind(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		user, err := finder.FindOne(c.Copy(), bson.M{"email": creds.Email})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("no user with the email exists: [email=%v]", creds.Email)})
			return
		}

		rawSessionId, exists := c.Get("sessionId")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("no sessionId is set in the context")})
			return
		}
		sessionId := rawSessionId.(string)

		sessions.Lock()
		session, err := sessions.FetchSession(c.Copy(), sessionId)
		if err != nil {
			sessions.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("sessionId no longer exists: [sessionId=%v]", sessionId)})
			return
		}
		session.UserId = user.ID.Hex()
		if err := sessions.UpdateSession(c.Copy(), sessionId, session); err != nil {
			sessions.Unlock()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sessions.Unlock()

		c.JSON(http.StatusOK, gin.H{"status": "logged in"})
	}
}

func RegisterAuthController(router *gin.RouterGroup, env *common.Env) {
	router.POST("/auth/login", Login(env.Repositories.UserRepository, env.Services.SessionService))
}
