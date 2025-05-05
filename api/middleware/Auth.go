package middleware

import (
	"go-auth-example/api/infra/cache"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type (
	Auth interface {
		SessionCookieAuth(c *gin.Context)
	}

	auth struct {
		redis cache.RedisClient
	}
)

const KeyUserID = "user_id"

func NewAuth(redis cache.RedisClient) Auth {
	return &auth{redis}
}

func (a *auth) SessionCookieAuth(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No session_id"})

		return
	}

	user_id, err := a.redis.Get(c, sessionID)
	if err != nil {
		if err == redis.Nil {
			// キー(sessionID)が存在しない場合もerrorが返ってくるらしい
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		}

		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.Set(KeyUserID, user_id) // 後続処理でc.Get(KeyUserID)で取得可能
	c.Next()
}
