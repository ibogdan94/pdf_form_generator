package middleware

//import (
//	"github.com/gin-gonic/gin"
//	"github.com/jinzhu/gorm"
//	"fmt"
//)

//func DBConnectHandler(db *gorm.DB) gin.HandlerFunc {
//func DBConnectHandler(string *string) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		_, ok := ctx.Set("db", string)
//
//		if !ok {
//			fmt.Println("Cannot set db connection to context")
//		}
//	}
//}
