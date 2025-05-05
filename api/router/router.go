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
	registUser controller.RegistUser,
	salogin controller.SessionLogin,
	salogout controller.SessionLogout,
) *gin.Engine {
	r := gin.Default()

	// ここでミドルウェアで認証したい
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			// ユーザー登録
			auth.POST("/user", registUser.Handle)

			// Basic認証
			auth.GET("/basic", ba.Handler)

			// sessin-cookie認証
			auth.POST("/session-cookie/login", salogin.Handler)
			auth.POST("/session-cookie/logout", salogout.Handler)
		}

		sessionAuth := api.Group("/session-auth")
		{
			// クッキー・セッション認証用の認証ミドルウェア
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
