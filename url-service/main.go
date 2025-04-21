package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"url-service/config"
	"url-service/internal/middleware"
	_urlMappingHandler "url-service/internal/url-mapping/delivery/http/handler"
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

	urlMappingRepository := _urlMappingPostgresDB.NewURLMappingRepo(config.PostgresDB)
	urlMappingUseCase := _urlMappingUseCase.NewURLMappingUseCase(urlMappingRepository)
	_urlMappingHandler.NewURLMappingHandler(ginEngine, urlMappingUseCase)

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
