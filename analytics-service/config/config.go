package config

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type Config struct {
	Port                string
	JWTSecret           []byte
	Rabbitmq            *amqp091.Channel
	URLClickedQueue     string
	ElasticSearchClient *elasticsearch.Client
	AnalyticsIndex      string
}

type ElasticSearchConfig struct {
	Host     string
	Port     string
	Username string
	Password string
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
	analyticsIndex := viper.GetString("ANALYTICS_INDEX")
	if analyticsIndex == "" {
		return nil, fmt.Errorf("missing ANALYTICS_INDEX")
	}
	urlClickedQueue := viper.GetString("URL_CLICKED_QUEUE")
	if urlClickedQueue == "" {
		return nil, fmt.Errorf("missing URL_CLICKED_QUEUE")
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

	elasticConfig := ElasticSearchConfig{
		Host:     viper.GetString("ELASTICSEARCH_HOST"),
		Port:     viper.GetString("ELASTICSEARCH_PORT"),
		Username: viper.GetString("ELASTICSEARCH_USERNAME"),
		Password: viper.GetString("ELASTICSEARCH_PASSWORD"),
	}

	ecnt, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%s", elasticConfig.Host, elasticConfig.Port)},
		Username:  elasticConfig.Username,
		Password:  elasticConfig.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	return &Config{
		Port:                port,
		JWTSecret:           []byte(jwtSecret),
		Rabbitmq:            ch,
		ElasticSearchClient: ecnt,
		URLClickedQueue:     urlClickedQueue,
		AnalyticsIndex:      analyticsIndex,
	}, nil

}
