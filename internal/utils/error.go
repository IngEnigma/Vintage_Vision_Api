package utils

import (
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, message string, err error) {
	Logger.Errorf("%s: %v", message, err)
	c.JSON(statusCode, gin.H{
		"error":    message,
		"detalles": err.Error(),
	})
}
