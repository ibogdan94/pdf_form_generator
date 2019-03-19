package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func hasErrorHandler(ctx *gin.Context, err error, message string, code int) bool {
	if err != nil {
		fmt.Printf("Error: %s\n", err)

		ctx.JSON(code, gin.H{
			"message": message,
		})

		return true
	}

	return false
}
