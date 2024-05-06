package services

import (
	"github.com/gin-gonic/gin"
)

type AuthenticateService interface {
	Login(c *gin.Context)
	GetUserInfo(c *gin.Context)
	GetNewToken(c *gin.Context)
	Logout(c *gin.Context)
	Register(c *gin.Context)
	ActiveAccount(c *gin.Context)
}
