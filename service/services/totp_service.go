package services

import "github.com/gin-gonic/gin"

type TotpService interface {
	GenerateQrCode(c *gin.Context)
	// GenerateTotp(c *gin.Context)
}
