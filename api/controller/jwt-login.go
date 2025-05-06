package controller

import (
	"fmt"
	"go-auth-example/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWTLogin interface {
	Handle(ctx *gin.Context)
}

type (
	jwtLogin struct {
		jwtissuer usecase.JWTIssuer
	}

	jwtLoginRequestPaarams struct {
		UserID   string
		Password string
	}
)

func NewJWTLogin(jwtissuer usecase.JWTIssuer) JWTLogin {
	return &jwtLogin{jwtissuer}
}

func (c *jwtLogin) Handle(ctx *gin.Context) {
	request := jwtLoginRequestPaarams{}
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("failed request bind", err))

		return
	}

	params := usecase.JWTIssuerParams{
		UserID:   request.UserID,
		Password: request.Password,
	}

	jwt, err := c.jwtissuer.Execute(params)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": jwt,
	})
}
