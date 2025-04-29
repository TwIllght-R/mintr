package config

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Port                 string
	JWTSecret            []byte
	RedirectExchange     string
	URLCreatedRoutingKey string
	URLUpdatedRoutingKey string
	URLDeletedRoutingKey string
	PostgresDB           *gorm.DB
	Rabbitmq             *amqp091.Channel
}

type PostgresDBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Rabbitmq struct {
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
	if port == "" {
		return nil, fmt.Errorf("missing port")
	}
	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("missing JWT secret")
	}
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("JWT secret must be at least 32 characters long")
	}
	redirectExchange := viper.GetString("REDIRECT_EXCHANGE")
	if redirectExchange == "" {
		return nil, fmt.Errorf("missing redirect exchange")
	}
	urlCreatedRoutingKey := viper.GetString("URL_CREATED_ROUTING_KEY")
	if urlCreatedRoutingKey == "" {
		return nil, fmt.Errorf("missing URL created routing key")
	}
	urlUpdatedRoutingKey := viper.GetString("URL_UPDATED_ROUTING_KEY")
	if urlUpdatedRoutingKey == "" {
		return nil, fmt.Errorf("missing URL updated routing key")
	}
	urlDeletedRoutingKey := viper.GetString("URL_DELETED_ROUTING_KEY")
	if urlDeletedRoutingKey == "" {
		return nil, fmt.Errorf("missing URL deleted routing key")
	}

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

	// Load AMQP config
	rabbitmqConfig := Rabbitmq{
		Host:     viper.GetString("RABBITMQ_HOST"),
		Port:     viper.GetString("RABBITMQ_PORT"),
		Username: viper.GetString("RABBITMQ_USERNAME"),
		Password: viper.GetString("RABBITMQ_PASSWORD"),
	}

	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqConfig.Username, rabbitmqConfig.Password, rabbitmqConfig.Host, rabbitmqConfig.Port))
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:                 port,
		JWTSecret:            []byte(jwtSecret),
		PostgresDB:           db,
		Rabbitmq:             ch,
		RedirectExchange:     redirectExchange,
		URLCreatedRoutingKey: urlCreatedRoutingKey,
		URLUpdatedRoutingKey: urlUpdatedRoutingKey,
		URLDeletedRoutingKey: urlDeletedRoutingKey,
	}, nil
}
