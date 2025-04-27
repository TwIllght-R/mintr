package main

import (
	"auth-service/config"
	"auth-service/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"time"

	_customerAuthHandler "auth-service/internal/customer-auth/delivery/http/handler"
	_notificationEventProducer "auth-service/internal/customer-auth/producer/rabbitmq"
	_customerAuthPostgresDB "auth-service/internal/customer-auth/repository/postgresdb"
	_customerAuthRedis "auth-service/internal/customer-auth/repository/redis"
	_customerAuthUseCase "auth-service/internal/customer-auth/usecase"
	_healthCheckHandler "auth-service/internal/healthcheck/http"
	_permissionPostgresDB "auth-service/internal/permission/repository/postgresdb"
	_permissionUseCase "auth-service/internal/permission/usecase"
	_tokenManager "auth-service/internal/token-manager"

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

	//permission manager
	permissionRepo := _permissionPostgresDB.NewPermissionRepository(config.PostgresDB)
	roleRepo := _permissionPostgresDB.NewRoleRepository(config.PostgresDB)
	rolePermissionRepo := _permissionPostgresDB.NewRolePermissionRepository(config.PostgresDB)
	permissionManager := _permissionUseCase.NewPermissionUsecase(permissionRepo, roleRepo, rolePermissionRepo)

	//token manager
	revocation := _tokenManager.NewRevocationStore(config.RedisDB)
	tokenManager := _tokenManager.NewTokenManager(config.JWTSecret, config.AccessTokenTTL, config.RefreshTokenTTL, revocation)

	customerAuthRepository := _customerAuthPostgresDB.NewCustomerAuthRepository(config.PostgresDB)
	customerAuthCache := _customerAuthRedis.NewCustomerAuthCache(config.RedisDB)
	notificationEventProducer := _notificationEventProducer.NewCustomerAuthEventProducer(config.Rabbitmq)

	customerAuthUseCase := _customerAuthUseCase.NewCustomerAuthUseCase(customerAuthRepository, customerAuthCache, notificationEventProducer, tokenManager, permissionManager)
	_customerAuthHandler.NewCustomerAuthHandler(ginEngine, customerAuthUseCase, config.JWTSecret)

	if err != nil {
		log.Println("Error signing up:", err)
	}
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
