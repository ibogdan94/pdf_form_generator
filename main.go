package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"os"
	"os/signal"
	"log"
	"context"
	"time"
	"net/http"
	"github.com/ibogdan94/pdf_form_generator/handlers"
	"github.com/ibogdan94/pdf_form_generator/utils"
)

func main() {
	props, err := utils.ParseJSONConfig()

	if err != nil {
		log.Fatal(err)
	}


	r := gin.Default()

	//gin.SetMode(gin.ReleaseMode)

	r.Use(gin.Logger())
	r.Use(location.Default())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("X-Requested-With", "Access-Control-Allow-Headers", "Cache-Control")

	r.Use(cors.New(config))

	r.Delims("{{", "}}")
	r.LoadHTMLFiles("./templates/home.html")
	r.Static("/static", "./static")
	r.Static("/node_modules", "./node_modules")
	r.Static("/temp", "./temp")

	r.GET("/", handlers.HomeHandler)
	r.POST("/api/v1/pdf/upload", handlers.ValidateUploadPDF)
	r.POST("/api/v1/pdf/save", handlers.SavePDF)

	srv := &http.Server{
		Addr: props.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServeTLS(props.ServerPem, props.ServerKey); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	//go func() {
	//	if err := srv.ListenAndServe(); err != nil {
	//		log.Printf("listen: %s\n", err)
	//	}
	//}()

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