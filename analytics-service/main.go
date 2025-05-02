package main

import (
	"analytics-service/config"
	_analyticsHandler "analytics-service/internal/analytics/delivery/http/handler"
	_eventHandler "analytics-service/internal/analytics/delivery/mq/handler"
	_analyticsElasticSearch "analytics-service/internal/analytics/repository/elasticsearch"
	_analyticsUseCase "analytics-service/internal/analytics/usecase"
	_healthCheckHandler "analytics-service/internal/healthcheck/http"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_analyticsRabbitMQ "analytics-service/internal/analytics/consumer/rabbitmq"
	"analytics-service/internal/middleware"
	"log"

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
	analyticsRepository := _analyticsElasticSearch.NewAnalyticsRepository(config.ElasticSearchClient, config.AnalyticsIndex)
	analyticsUseCase := _analyticsUseCase.NewAnalyticsUseCase(analyticsRepository)
	_analyticsHandler.NewAnalyticsHandler(ginEngine, config.JWTSecret, analyticsUseCase)

	urlClickedEventHandler := _eventHandler.NewURLClickedEventHandler(analyticsUseCase)

	urlClickedConsumer := _analyticsRabbitMQ.NewAnalyticsConsumer(config.Rabbitmq, config.URLClickedQueue, urlClickedEventHandler)

	go func() {
		if err := urlClickedConsumer.Start(); err != nil {
			log.Fatalf("Failed to start visit consumer: %v", err)
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

	if err := urlClickedConsumer.Close(); err != nil {
		log.Fatalf("Failed to close visit consumer: %v", err)
	}

	log.Println("🛑 Server gracefully stopped")

}
