package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pdf_form_generator/handlers"
	"time"
)

type baseConfig struct {
	pwd string
}

var serverConfig baseConfig

func init() {
	pwd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Cannot get pwd: %s\n", err)
	}

	serverConfig.pwd = pwd
}

func main() {
	r := gin.Default()

	r.Use(location.Default())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("X-Requested-With", "Access-Control-Allow-Headers", "Cache-Control")

	r.Use(cors.New(config))

	h := handlers.NewHandlers(serverConfig.pwd)
	h.SetupRoutes(r)

	srv := &http.Server{
		Addr:    ":8080",
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
