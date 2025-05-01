package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-service/config"
	_healthCheckHandler "url-service/internal/healthcheck/http"
	"url-service/internal/middleware"
	_urlMappingHandler "url-service/internal/url-mapping/delivery/http/handler"
	_urlMappingRabbitMQ "url-service/internal/url-mapping/producer/rabbitmq"
	_urlMappingPostgresDB "url-service/internal/url-mapping/repository/postgresdb"
	_urlMappingUseCase "url-service/internal/url-mapping/usecase"

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

	urlMappingRepository := _urlMappingPostgresDB.NewURLMappingRepository(config.PostgresDB)
	urlMappingProducer := _urlMappingRabbitMQ.NewURLMappingEventProducer(config.Rabbitmq, config.RedirectExchange, config.URLCreatedRoutingKey, config.URLUpdatedRoutingKey, config.URLDeletedRoutingKey)
	urlMappingUseCase := _urlMappingUseCase.NewURLMappingUseCase(urlMappingRepository, urlMappingProducer)
	_urlMappingHandler.NewURLMappingHandler(ginEngine, config.JWTSecret, urlMappingUseCase)

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

	sqlDB, err := config.PostgresDB.DB()
	if err != nil {
		log.Printf("Failed to get SQL DB from GORM: %v", err)
	} else if err := sqlDB.Close(); err != nil {
		log.Printf("Failed to close PostgresDB client: %v", err)
	}

	if err := config.Rabbitmq.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ client: %v", err)
	}

	log.Println("🛑 Server gracefully stopped")
}
