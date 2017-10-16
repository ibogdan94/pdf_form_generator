package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func HomeHandler(c *gin.Context)  {
	fmt.Println(c.Get("db"))
	c.HTML(http.StatusOK, "home.html", map[string]interface{}{
		"title": "Home Page",
	})
}