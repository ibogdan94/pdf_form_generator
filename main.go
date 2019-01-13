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

//func init() {
//	_, err := utils.ParseJSONConfig()
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}

func init() {
	pwd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Cannot get pwd: %s\n", err)
	}

	serverConfig.pwd = pwd
}

func main() {
	//props := utils.Config

	r := gin.Default()

	//if props.Env == "prod" {
	//	gin.SetMode(gin.ReleaseMode)
	//	gin.DisableConsoleColor()
	//} else {
	//	gin.SetMode(gin.DebugMode)
	//}
	//
	//pwd, err := os.Getwd()
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//resultFolder := pwd + "/" + utils.Config.TempPath + "/result"
	//
	//if _, err := os.Stat(resultFolder); os.IsNotExist(err) {
	//	if err := os.Mkdir(resultFolder, 0755); err != nil {
	//		log.Fatalf("Cannot create result folder %s. Error: %s", resultFolder, err)
	//	}
	//}
	//
	//f, err := os.Create(pwd + "/logs/" + time.Now().Format("20060102150405") + ".log")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	//log.SetOutput(gin.DefaultWriter)
	//
	r.Use(location.Default())
	//
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("X-Requested-With", "Access-Control-Allow-Headers", "Cache-Control")
	//
	r.Use(cors.New(config))
	//
	//r.Delims("{{", "}}")
	//r.LoadHTMLFiles("./templates/home.html")
	//r.Static("/static", "./static")
	//r.Static("/node_modules", "./node_modules")
	//r.Static("/temp", "./temp")
	//
	//r.GET("/", handlers.HomeHandler)
	//r.POST("/api/v1/pdf/upload", handlers.ValidateUploadPDF)
	//r.GET("/api/v1/pdf/:code", handlers.GetImagesByCode)
	//r.POST("/api/v1/pdf/generate/:code", handlers.GeneratePDF)
	//r.POST("/api/v1/pdf/cleanup", handlers.Cleanup)

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
