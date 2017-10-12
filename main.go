package main

import (
	"github.com/gin-gonic/gin"
	handlers "./handlres"
	"os"
	"os/signal"
	"log"
	"context"
	"time"
	"net/http"
)

const WEB_SERVER_PORT = ":8888"

func main() {
	r := gin.Default()

	r.Delims("{{", "}}")
	r.LoadHTMLFiles("./templates/home/home.html", "./templates/builder/builder.html")
	r.Static("/public", "./public")
	r.Static("/node_modules", "./node_modules")
	r.Static("/temp", "./temp")

	r.GET("/", handlers.HomeHandler)
	r.POST("/pdf/upload", handlers.ValidateUploadPDF)
	r.GET("/pdf/edit", handlers.EditParsedPDF)

	srv := &http.Server{
		Addr: WEB_SERVER_PORT,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
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
