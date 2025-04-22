package main

import (
	"fmt"
	"log"
	"net/http"
	"redirect-service/config"
	_healthCheckHandler "redirect-service/internal/healthcheck/http"
	"redirect-service/internal/middleware"
	_redirectURLHandler "redirect-service/internal/redirect-url/delivery/http/handler"
	_redirectURLPostgresDB "redirect-service/internal/redirect-url/repository/postgresdb"
	_redirectURLCache "redirect-service/internal/redirect-url/repository/redis"
	_redirectURLUseCase "redirect-service/internal/redirect-url/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	config, err := config.New()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	ginEngine := gin.New()

	ginEngine.Use(
		middleware.Logger(),
		gin.Recovery(),
		middleware.RequestID(),
		middleware.CORS(),
	)

	//health check
	_healthCheckHandler.NewHealthCheckHandler(ginEngine)

	redirectURLRepository := _redirectURLPostgresDB.NewRedirectURLRepository(config.PostgresDB)
	redirectCache := _redirectURLCache.NewRedirectURLCache(config.RedisDB)
	redirectURLUseCase := _redirectURLUseCase.NewRedirectURLUsecase(redirectURLRepository, redirectCache)
	_redirectURLHandler.NewRedirectURLHandler(ginEngine, redirectURLUseCase)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.Port),
		Handler:        ginEngine,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	log.Printf("Starting server on port %s", config.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
