package services

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthenticateHandler struct {
	DB *gorm.DB
}

func (h *AuthenticateHandler) Login(c *gin.Context) {
}
