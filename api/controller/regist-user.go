package controller

import (
	"go-auth-example/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistUser interface {
	Handle(ctx *gin.Context)
}

type (
	registUser struct {
		register usecase.UserRegister
	}

	registUserRequest struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}
)

func NewRegistUser(register usecase.UserRegister) RegistUser {
	return &registUser{register}
}

func (c *registUser) Handle(ctx *gin.Context) {
	request := registUserRequest{}
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	params := usecase.UserRegisterParams{
		UserID:   request.UserID,
		Password: request.Password,
	}

	if err := c.register.Execute(ctx, params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	ctx.Status(http.StatusCreated)
}
