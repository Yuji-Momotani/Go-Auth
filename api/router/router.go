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
	jlogin controller.JWTLogin,
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
			// auth.POST("/session-cookie/logout", salogout.Handler)
			auth.POST("/jwt/login", jlogin.Handle)
			// auth.POST("/session-cookie/logout", salogout.Handler)
		}

		sessionCookie := api.Group("/session-cookie")
		{
			// クッキー・セッション認証用の認証ミドルウェア
			sessionCookie.Use(m.SessionCookieAuth)
			// API認証ができているかチェックするためのハンドラー
			sessionCookie.GET("/hello", helloHandle)
		}

		jwt := api.Group("/jwt")
		{
			// JWT用の認証ミドルウェア
			jwt.Use(m.JWTAuth)
			// API認証ができているかチェックするためのハンドラー
			jwt.GET("/hello", helloHandle)
		}

	}

	return r
}

func helloHandle(ctx *gin.Context) {
	userID, exists := ctx.Get(middleware.KeyUserID)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})

		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("hello %s", userID)})
}
