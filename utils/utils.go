package utils

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func GetSchemaAndHost(ctx *gin.Context) (schemaAndHost string) {
	url := location.Get(ctx)

	schemaAndHost = url.Scheme + "://" + url.Host
	return
}