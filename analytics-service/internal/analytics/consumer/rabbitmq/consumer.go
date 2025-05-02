package rabbitmq

import (
	"analytics-service/domain/interfaces"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type analyticsConsumer struct {
	ch      *amqp091.Channel
	queue   string
	handler interfaces.EventHandler
}

func NewAnalyticsConsumer(ch *amqp091.Channel, queue string, handler interfaces.EventHandler) *analyticsConsumer {
	return &analyticsConsumer{
		ch:      ch,
		queue:   queue,
		handler: handler,
	}
}

func (c *analyticsConsumer) Start() error {
	msgs, err := c.ch.Consume(c.queue, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", msg.Body)
			c.handler.HandleEvent(msg.Body)
		}
	}()

	log.Printf("RabbitMQ consumer started for queue: %s", c.queue)
	return nil
}

func (c *analyticsConsumer) Close() error {
	if err := c.ch.Close(); err != nil {
		return err
	}
	log.Println("RabbitMQ consumer closed.")
	return nil
}
