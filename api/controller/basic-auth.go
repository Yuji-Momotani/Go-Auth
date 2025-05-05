package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	BasicAuth interface {
		Handler(c *gin.Context)
	}

	basicAuth struct{}
)

func NewBasicAuth() BasicAuth {
	return &basicAuth{}
}

func (h *basicAuth) Handler(c *gin.Context) {
	const (
		USERNAME = "admin"
		PASSWORD = "password123"
	)

	user, pass, ok := c.Request.BasicAuth()

	if !ok {
		resp := gin.H{
			"message:": "BadRequest",
			"code":     "INVALID_REQUEST",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if user == USERNAME && pass == PASSWORD {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome, admin!"})
		return
	}

	c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
	c.Status(http.StatusUnauthorized)
}
