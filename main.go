package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"log"
	"context"
	"time"
	"net/http"
	"github.com/jinzhu/gorm"
	 _ "github.com/jinzhu/gorm/dialects/mysql"
	handlers "./handlers"
	"io"
)

const WEB_SERVER_PORT = ":8888"

func DBConnectHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("DB", db)
	}
}

func main() {
	db, err := gorm.Open("mysql", "ddidev:ddipass@tcp(127.0.0.1:3307)/akela?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Print(err)
	}

	defer db.Close()

	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()

	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(DBConnectHandler(db))

	r.Delims("{{", "}}")
	r.LoadHTMLFiles("./templates/home/home.html", "./templates/builder/builder.html")
	r.Static("/static", "./static")
	r.Static("/node_modules", "./node_modules")
	r.Static("/temp", "./temp")

	r.GET("/", handlers.HomeHandler)
	r.POST("/pdf/upload", handlers.ValidateUploadPDF)
	r.GET("/pdf/edit", handlers.EditParsedPDF)
	r.POST("/pdf/save", handlers.SavePDF)

	srv := &http.Server{
		Addr: WEB_SERVER_PORT,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()


	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exist")
}