package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sochain-client/pkg/controller"
	"sochain-client/pkg/sochain"
	"sochain-client/pkg/util"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("missing .env file")
	}

	logger, err := util.NewLogger(util.GetEnv("CI_ENV", util.EnvironmentDev))
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	client := sochain.NewSochain()
	controller := controller.NewController(logger, client)
	RegisterRoutes(r, controller)

	srv := &http.Server{
		Addr:    util.GetEnv("HOST", "localhost") +":"+ util.GetEnv("API_PORT", "8080"), 
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("server shutdown started...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown was forced %v", err)
	}
}

func RegisterRoutes(e *gin.Engine, c *controller.Controller) {
	e.GET("/network/:id", c.HandleGetBlock)
	e.GET("/network/:id/tx/:txhash", c.HandleGetTransaction)
}
