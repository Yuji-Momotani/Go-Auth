package middleware

import (
	"fmt"
	"go-auth-example/api/infra/cache"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type (
	Auth interface {
		SessionCookieAuth(c *gin.Context)
		JWTAuth(c *gin.Context)
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

func (a *auth) JWTAuth(c *gin.Context) {
	// 1. Authorization ヘッダー取得
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "missing authorization",
		})

		return
	}

	// 2. Bearer スキームチェック
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authorization header format",
		})
	}

	tokenString := parts[1]

	// 3. JWTパース&検証
	secret := []byte(os.Getenv("JWT_SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// アルゴリズムのチェック
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("expected signing method: %v", t.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})

		return
	}

	// 4. ユーザーIDの取り出し
	clams, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		// 有効期限のチェックはライブラリの中で行われており、token.Validで判定できる
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
	}

	sub, ok := clams["sub"].(string)
	if !ok || sub == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
	}
	c.Set(KeyUserID, sub)
	c.Next()
}
