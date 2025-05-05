package router

import (
	"fmt"
	"go-auth-example/api/controller"
	"go-auth-example/api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	m middleware.Auth,
	ba controller.BasicAuth,
	salogin controller.SessionLogin,
	salogout controller.SessionLogout,
) *gin.Engine {
	r := gin.Default()

	// Basic認証ルート
	r.GET("/basic-auth", ba.Handler)

	// sessin-cookie認証
	r.POST("/session-auth/login", salogin.Handler)
	r.POST("/session-auth/logout", salogout.Handler)

	// ここでミドルウェアで認証したい
	api := r.Group("/api")
	{
		sessionAuth := api.Group("/session-auth")
		{
			sessionAuth.Use(m.SessionCookieAuth)
			sessionAuth.GET("/hello", func(ctx *gin.Context) {
				userID, exists := ctx.Get(middleware.KeyUserID)
				if !exists {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})

					return
				}
				ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("hello %s", userID)})
			})
		}
	}

	return r
}
