package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pdf_form_generator/handlers"
	"time"
)

type baseConfig struct {
	pwd     string
	Port    int    `json:"port"`
	Mode    string `json:"mode"`
	BaseURL string `json:"baseURL"`
}

var config baseConfig

func init() {
	pwd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Cannot get pwd: %s\n", err)
	}

	config.pwd = pwd

	payload, err := ioutil.ReadFile(pwd + "/config.json")

	if err != nil {
		log.Fatalf("Something went wrong with config.json: %s\n", err)
	}

	if err := json.Unmarshal(payload, &config); err != nil {
		log.Fatalf("Cannot unmarshal config.json to go struct: %s\n", err)
	}
}

func main() {
	gin.SetMode(config.Mode)
	r := gin.Default()

	r.Static("/result", "./result")
	r.Static("/store", "./store")

	h := handlers.NewHandlers(config.pwd, config.BaseURL)
	h.SetupRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
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
