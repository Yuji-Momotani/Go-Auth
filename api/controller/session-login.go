package controller

import (
	"go-auth-example/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	SessionLogin interface {
		Handler(c *gin.Context)
	}

	sessionLogin struct {
		usecase usecase.SessionLogin
	}

	Request struct {
		UserID   string `json:"user_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

func NewSessionLogin(usecase usecase.SessionLogin) SessionLogin {
	return &sessionLogin{usecase}
}

func (h *sessionLogin) Handler(c *gin.Context) {
	var req Request
	if err := c.Bind(&req); err != nil {
		c.Status(http.StatusBadRequest)

		return
	}

	params := usecase.SessionLoginParams{
		UserID:   req.UserID,
		Password: req.Password,
	}

	sessionID, err := h.usecase.Execute(c, params)
	if err != nil {
		c.Status(http.StatusInternalServerError)

		return
	}

	// if req.UserName != USERNAME || req.Password != PASSWORD {
	// 	c.Status(http.StatusUnauthorized)
	// 	return
	// }

	// session_id, _ := uuid.NewUUID()
	c.SetCookie("session_id", sessionID, 3600, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}
