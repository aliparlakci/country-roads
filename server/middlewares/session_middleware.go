package middlewares

import (
	"fmt"
	"github.com/aliparlakci/country-roads/services"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware(repo services.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sessionId string
		var session services.Session
		id, err := c.Cookie("visitor")
		if err != nil { // no session cookie exists
			session = services.Session{}
			if sessionId, err = repo.CreateSession(c.Copy(), session); err != nil {
				panic(err)
			}
			c.Header("Set-Cookie", fmt.Sprintf("visitor=%v; Max-Age=7776000;", sessionId))
		} else { // session cookie exists
			if session, err = repo.FetchSession(c.Copy(), id); err != nil { // given session cookie is expired
				session = services.Session{}
				if sessionId, err = repo.CreateSession(c.Copy(), session); err != nil {
					panic(err)
				}
				c.Header("Set-Cookie", fmt.Sprintf("visitor=%v; Max-Age=7776000;", sessionId))
			}
		}
		c.Set("session", session)
		c.Next()
	}
}
