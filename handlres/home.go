package handlres

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeHandler(c *gin.Context)  {
	c.HTML(http.StatusOK, "home.html", map[string]interface{}{
		"title": "Home Page",
	})
}