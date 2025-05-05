package controller

import (
	"github.com/gin-gonic/gin"
)

type (
	SessionLogout interface {
		Handler(c *gin.Context)
	}

	sessionLogout struct{}
)

func NewSessionLogout() SessionLogout {
	return &sessionLogout{}
}

func (h *sessionLogout) Handler(c *gin.Context) {
	// cookieの削除
	c.SetCookie("session_id", "", -1, "/", "localhost", false, true)
}
