package controllers

import (
	"fmt"
	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Login(auth services.AuthService, finder models.UserFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var creds models.LoginRequestForm
		if err := c.Bind(&creds); err != nil {
			logger.WithField("error", err.Error()).Debug("cannot bind request body to models.LoginRequestForm")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		user, err := finder.FindOne(c.Copy(), bson.M{"email": creds.Email})
		if err != nil {
			logger.WithField("email", creds.Email).Debug("no user with email exists")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("no user with the email exists: [email=%v]", creds.Email)})
			return
		}

		if err := auth.CreateOTP(user.ID.Hex()); err != nil {
			logger.WithField("error", err.Error()).Error("auth.CreateOTP raised an error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		logger.WithFields(logrus.Fields{
			"user_id": user.ID.Hex(),
			"email": user.Email,
		}).Info("login flow is started for user with user_id and email")
		c.JSON(http.StatusOK, gin.H{})
	}
}

func Verify(auth services.AuthService, sessions services.SessionService, userFinder models.UserFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var creds models.VerifyRequestForm
		if err := c.Bind(&creds); err != nil {
			logger.WithField("error", err).Debug("cannot bind request body to models.VerifyRequestForm")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		rawSessionId, exists := c.Get("sessionId")
		if !exists {
			logger.Error("cannot find sessionId parameter in the context")
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("no sessionId is set in the context")})
			return
		}

		user, err := userFinder.FindOne(c.Copy(), bson.M{"email": creds.Email})
		if err != nil {
			logger.WithField("email", creds.Email).Debug("no user exists with the email")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("no user with the email exists: [email=%v]", creds.Email)})
			return
		}

		verified, err := auth.VerifyOTP(user.ID.Hex(), creds.OTP)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user_id": user.ID.Hex(),
				"email": creds.Email,
			}).Debug("no OTP has found for user with user_id and email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
			return
		}
		if !verified {
			logger.WithFields(logrus.Fields{
				"otp": creds.OTP,
				"user_id": user.ID.Hex(),
				"email": creds.Email,
			}).Debug("otp does not match for user with user_id and email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP did not match"})
			return
		}

		sessionId := rawSessionId.(string)
		sessions.Lock()
		session, err := sessions.FetchSession(c.Copy(), sessionId)
		if err != nil {
			sessions.Unlock()
			logger.WithField("id", sessionId).Debug("session with id no longer exists")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("sessionId no longer exists: [sessionId=%v]", sessionId)})
			return
		}
		session.UserId = user.ID.Hex()
		if err := sessions.UpdateSession(c.Copy(), sessionId, session); err != nil {
			sessions.Unlock()
			logger.WithFields(logrus.Fields{
				"id": sessionId,
				"error": err.Error(),
			}).Error("session with sessionId cannot get updated, SessionService.UpdateSession() raised an error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sessions.Unlock()

		auth.RevokeOTP(user.ID.Hex())
		logger.WithFields(logrus.Fields{
			"id": user.ID.Hex(),
			"email": user.Email,
		}).Info("OTP for user with id and email is revoked")

		logger.WithFields(logrus.Fields{
			"user_id": user.ID.Hex(),
			"email": user.Email,
			"sessionId": sessionId,
		}).Info("user with user_id and email is logged in in the session with sessionId")
		c.JSON(http.StatusOK, gin.H{"status": "logged in"})
	}
}

func RegisterAuthController(router *gin.RouterGroup, env *common.Env) {
	router.POST("/auth/login", Login(
		env.Services.AuthService,
		env.Repositories.UserRepository,
	))
	router.POST("/auth/verify", Verify(
		env.Services.AuthService,
		env.Services.SessionService,
		env.Repositories.UserRepository,
	))
}
