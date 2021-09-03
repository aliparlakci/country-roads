package controllers

import (
	"fmt"
	"github.com/aliparlakci/country-roads/common"
	"github.com/aliparlakci/country-roads/models"
	"github.com/aliparlakci/country-roads/repositories"
	"github.com/aliparlakci/country-roads/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Login(otpService services.OTPService, finder repositories.UserFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var creds models.LoginRequestForm
		if err := c.Bind(&creds); err != nil {
			logger.Debugf("cannot bind request body to models.LoginRequestForm: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		if user, isLoggedIn := c.Get("user"); isLoggedIn {
			logger.WithField("user_id", user.(models.User).ID).Debug("user with user_id is already logged in on this session")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is already logged in"})
			return
		}

		user, err := finder.FindOne(c.Copy(), bson.M{"email": creds.Email})
		if err != nil {
			logger.WithField("email", creds.Email).Debug("no user with email exists")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("no user with the email exists: [email=%v]", creds.Email)})
			return
		}

		if err := otpService.CreateOTP(user.ID.Hex()); err != nil {
			logger.WithField("user_id", user.ID.Hex()).Errorf("while creating OTP for user with user_id, otpService.CreateOTP raised an error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		logger.WithFields(logrus.Fields{
			"user_id": user.ID.Hex(),
			"email":   user.Email,
		}).Debug("login flow is started for user with user_id and email")
		c.JSON(http.StatusOK, gin.H{})
	}
}

func Verify(otpService services.OTPService, sessions services.SessionService, users repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var creds models.VerifyRequestForm
		if err := c.Bind(&creds); err != nil {
			logger.Debugf("cannot bind request body to models.VerifyRequestForm: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		user, err := users.FindOne(c.Copy(), bson.M{"email": creds.Email})
		if err != nil {
			logger.WithField("email", creds.Email).Debug("no user exists with the email")
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("no user with the email exists: [email=%v]", creds.Email)})
			return
		}

		correct, err := otpService.VerifyOTP(user.ID.Hex(), creds.OTP)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user_id": user.ID.Hex(),
				"email":   creds.Email,
			}).Debug("no OTP has found for user with user_id and email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
			return
		}
		if !correct {
			logger.WithFields(logrus.Fields{
				"otp":     creds.OTP,
				"user_id": user.ID.Hex(),
				"email":   creds.Email,
			}).Debug("otp does not match for user with user_id and email")
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP did not match"})
			return
		}

		sessionId, err := sessions.CreateSession(c.Copy(), user.ID.Hex())
		if err != nil {
			logger.WithField("user_id", user.ID.Hex()).Errorf("SessionService.CreateSession raised an error while creating a new session for user with user_id: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot log in"})
			return
		}
		c.Header("Set-Cookie", fmt.Sprintf("session=%v; Max-Age=7776000; path=/;", sessionId))

		if err := otpService.RevokeOTP(user.ID.Hex()); err != nil {
			logger.WithField("user_id", user.ID.Hex()).Warn("cannot revoke otp for user with user_id: %v", err.Error())
		} else {
			logger.WithFields(logrus.Fields{
				"id":    user.ID.Hex(),
				"email": user.Email,
			}).Info("OTP for user with id and email is revoked")
		}

		if !user.Verified {
			if err := users.UpdateOne(
				c.Copy(),
				bson.M{"_id": user.ID},
				bson.D{{"$set", bson.D{{"verified", true}}}},
			); err != nil {
				logger.WithField("user_id", user.ID.Hex()).Errorf("RideRepository.UpdateOne raised an error while setting user verified: %v", err)
			}
		}

		logger.WithFields(logrus.Fields{
			"user_id": user.ID.Hex(),
			"email":   user.Email,
		}).Info("user with user_id and email is logged in")
		c.JSON(http.StatusOK, gin.H{"status": "logged in"})
	}
}

func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())

		var user models.User
		if value, exists := c.Get("user"); !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
			return
		} else {
			var success bool
			if user, success = value.(models.User); !success {
				logger.Errorf("type assertion for user in context failed")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot find user"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func Logout(sessions services.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c)
		sessionId, err := c.Cookie("session")
		if err != nil {
			logger.Debug("no user is logged in")
			c.JSON(http.StatusBadRequest, gin.H{"error": "no user is logged in"})
			return
		}

		userId, err := sessions.FetchSession(c.Copy(), sessionId)
		if err != nil {
			logger.WithField("session_id", sessionId).Info("session with session_id no longer exists")
			c.Header("Set-Cookie", fmt.Sprintf("session=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;"))
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		if err := sessions.RevokeSession(c.Copy(), sessionId); err != nil {
			logger.Errorf("cannot revoke session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot log you out"})
			return
		}
		c.Header("Set-Cookie", fmt.Sprintf("session=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/;"))

		logger.WithFields(logrus.Fields{
			"user_id":    userId,
			"session_id": sessionId,
		}).Info("user with user_id is logged out of the session with session_id")

		c.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterAuthController(router *gin.RouterGroup, env *common.Env) {
	router.POST("/auth/login", Login(
		env.Services.OTPService,
		env.Repositories.UserRepository,
	))
	router.POST("/auth/verify", Verify(
		env.Services.OTPService,
		env.Services.SessionService,
		env.Repositories.UserRepository,
	))
	router.GET("/auth/user", User())
	router.POST("/auth/logout", Logout(
		env.Services.SessionService,
	))
}
