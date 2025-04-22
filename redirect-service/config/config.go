package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Port       string
	PostgresDB *gorm.DB
	RedisDB    *redis.Client
}

type PostgresDBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func New() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("No .env file, using environment variables")
	}

	viper.AutomaticEnv()

	// Load server port
	port := viper.GetString("PORT")

	// Load Postgres config
	pgConfig := PostgresDBConfig{
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		Username: viper.GetString("POSTGRES_USERNAME"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("POSTGRES_DBNAME"),
		SSLMode:  viper.GetString("POSTGRES_SSLMODE"),
	}

	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		pgConfig.Host, pgConfig.Port, pgConfig.Username, pgConfig.DBName, pgConfig.Password, pgConfig.SSLMode)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	//Load Redis config
	rdbConfig := RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Username: viper.GetString("REDIS_USERNAME"),
		Password: viper.GetString("REDIS_PASSWORD"),
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     rdbConfig.Host + ":" + rdbConfig.Port,
		Username: rdbConfig.Username,
		Password: rdbConfig.Password,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &Config{
		Port:       port,
		PostgresDB: db,
		RedisDB:    rdb,
	}, nil
}
