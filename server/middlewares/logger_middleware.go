package middlewares

import (
	"github.com/aliparlakci/country-roads/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := common.LoggerWithRequestId(c.Copy())
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		uri := path
		if raw != "" {
			uri = uri + "?" + raw
		}

		/*requestFields := logrus.Fields{
			"time_stamp":  start.Format("2006/01/02 - 15:04:05"),
			"method":      c.Request.Method,
			"path":        uri,
		}*/

		// logger.WithFields(requestFields).Debug("request is received")

		// Process request
		c.Next()

		now := time.Now()

		responseFields := logrus.Fields{
			"time_stamp":  now.Format("2006/01/02 - 15:04:05"),
			"status_code": c.Writer.Status(),
			"latency":     now.Sub(start),
			"method":      c.Request.Method,
			"path":        uri,
		}

		if error := c.Errors.ByType(gin.ErrorTypePrivate).String(); error != "" {
			responseFields["error"] = error
		}

		logger.WithFields(responseFields).Debug()
	}
}