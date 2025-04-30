package config

import (
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Config struct {
	Port                  string
	RedisClient           *redis.Client
	MongoDB               *mongo.Database
	RedirectURLCollection string
	Rabbitmq              *amqp091.Channel
	URLCreatedQueue       string
	URLUpdatedQueue       string
	URLDeletedQueue       string
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type MongoDBConfig struct {
	Host          string
	Port          string
	Username      string
	Password      string
	DBName        string
	AuthSource    string
	AuthMechanism string
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
	redirectURLCollection := viper.GetString("REDIRECT_URL_COLLECTION")
	if redirectURLCollection == "" {
		return nil, fmt.Errorf("missing redirect URL collection name")
	}
	urlCreatedQueue := viper.GetString("URL_CREATED_QUEUE")
	if urlCreatedQueue == "" {
		return nil, fmt.Errorf("missing URL created queue")
	}
	urlUpdatedQueue := viper.GetString("URL_UPDATED_QUEUE")
	if urlUpdatedQueue == "" {
		return nil, fmt.Errorf("missing URL updated queue")
	}
	urlDeletedQueue := viper.GetString("URL_DELETED_QUEUE")
	if urlDeletedQueue == "" {
		return nil, fmt.Errorf("missing URL deleted queue")
	}

	// Load Mongo config
	mongoConfig := MongoDBConfig{
		Host:          viper.GetString("MONGO_HOST"),
		Port:          viper.GetString("MONGO_PORT"),
		Username:      viper.GetString("MONGO_USERNAME"),
		Password:      viper.GetString("MONGO_PASSWORD"),
		DBName:        viper.GetString("MONGO_DB_NAME"),
		AuthSource:    viper.GetString("MONGO_AUTH_SOURCE"),
		AuthMechanism: viper.GetString("MONGO_AUTH_MECHANISM"),
	}

	mcnt, err := mongo.Connect(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s&authMechanism=%s",
		mongoConfig.Username, mongoConfig.Password, mongoConfig.Host, mongoConfig.Port, mongoConfig.DBName, mongoConfig.AuthSource, mongoConfig.AuthMechanism)))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	mdb := mcnt.Database(mongoConfig.DBName)
	if err := mdb.Client().Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// //Load Redis config
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

	// Load AMQP config
	rabbitmqConfig := Rabbitmq{
		Host:     viper.GetString("RABBITMQ_HOST"),
		Port:     viper.GetString("RABBITMQ_PORT"),
		Username: viper.GetString("RABBITMQ_USERNAME"),
		Password: viper.GetString("RABBITMQ_PASSWORD"),
	}

	conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqConfig.Username, rabbitmqConfig.Password, rabbitmqConfig.Host, rabbitmqConfig.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return &Config{
		Port:                  port,
		MongoDB:               mdb,
		RedisClient:           rdb,
		Rabbitmq:              ch,
		URLCreatedQueue:       urlCreatedQueue,
		URLUpdatedQueue:       urlUpdatedQueue,
		URLDeletedQueue:       urlDeletedQueue,
		RedirectURLCollection: redirectURLCollection,
	}, nil
}
