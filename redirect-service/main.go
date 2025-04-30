package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"redirect-service/config"
	_healthCheckHandler "redirect-service/internal/healthcheck/http"
	"redirect-service/internal/middleware"
	"syscall"

	_redirectURLRabbitMQ "redirect-service/internal/redirect-url/consumer/rabbitmq"
	_redirectURLHandler "redirect-service/internal/redirect-url/delivery/http/handler"
	_eventHandler "redirect-service/internal/redirect-url/delivery/mq/handler"
	_redirectURLMongoDB "redirect-service/internal/redirect-url/repository/mongodb"
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
	redirectURLRepository := _redirectURLMongoDB.NewRedirectURLRepository(config.MongoDB, config.RedirectURLCollection)
	redirectCache := _redirectURLCache.NewRedirectURLCache(config.RedisClient)
	redirectURLUseCase := _redirectURLUseCase.NewRedirectURLUseCase(redirectURLRepository, redirectCache)
	_redirectURLHandler.NewRedirectURLHandler(ginEngine, redirectURLUseCase)

	urlCreatedEventHandler := _eventHandler.NewURLCreatedEventHandler(redirectURLUseCase)
	urlUpdatedEventHandler := _eventHandler.NewURLUpdatedEventHandler(redirectURLUseCase)
	urlDeletedEventHandler := _eventHandler.NewURLDeletedEventHandler(redirectURLUseCase)

	urlCreatedConsumer := _redirectURLRabbitMQ.NewRedirectURLConsumer(config.Rabbitmq, config.URLCreatedQueue, urlCreatedEventHandler)
	urlUpdatedConsumer := _redirectURLRabbitMQ.NewRedirectURLConsumer(config.Rabbitmq, config.URLUpdatedQueue, urlUpdatedEventHandler)
	urlDeletedConsumer := _redirectURLRabbitMQ.NewRedirectURLConsumer(config.Rabbitmq, config.URLDeletedQueue, urlDeletedEventHandler)

	go func() {
		if err := urlCreatedConsumer.Start(); err != nil {
			log.Fatalf("Failed to start created consumer: %v", err)
		}
	}()
	go func() {
		if err := urlUpdatedConsumer.Start(); err != nil {
			log.Fatalf("Failed to start updated consumer: %v", err)
		}
	}()
	go func() {
		if err := urlDeletedConsumer.Start(); err != nil {
			log.Fatalf("Failed to start deleted consumer: %v", err)
		}
	}()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.Port),
		Handler:        ginEngine,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		log.Printf("🟢 Starting server on port %s", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := urlCreatedConsumer.Close(); err != nil {
		log.Printf("Failed to close created consumer: %v", err)
	}
	if err := urlUpdatedConsumer.Close(); err != nil {
		log.Printf("Failed to close updated consumer: %v", err)
	}
	if err := urlDeletedConsumer.Close(); err != nil {
		log.Printf("Failed to close deleted consumer: %v", err)
	}
	if err := config.RedisClient.Close(); err != nil {
		log.Printf("Failed to close Redis client: %v", err)
	}
	if err := config.MongoDB.Client().Disconnect(context.Background()); err != nil {
		log.Printf("Failed to disconnect MongoDB client: %v", err)
	}
	if err := config.Rabbitmq.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ client: %v", err)
	}

	log.Println("🛑 Server gracefully stopped")

}
