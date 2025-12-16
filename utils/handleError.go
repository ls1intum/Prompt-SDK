package utils

import "github.com/gin-gonic/gin"

//nolint:unused // Public SDK function for external use
func handleError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
