package rabbitmq

import (
	"log"
	"redirect-service/domain/interfaces"

	"github.com/rabbitmq/amqp091-go"
)

type redirectURLConsumer struct {
	ch      *amqp091.Channel
	queue   string
	handler interfaces.EventHandler
}

func NewRedirectURLConsumer(ch *amqp091.Channel, queue string, handler interfaces.EventHandler) *redirectURLConsumer {
	return &redirectURLConsumer{
		ch:      ch,
		queue:   queue,
		handler: handler,
	}
}

func (c *redirectURLConsumer) Start() error {
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

func (c *redirectURLConsumer) Close() error {
	if err := c.ch.Close(); err != nil {
		return err
	}
	log.Println("RabbitMQ consumer closed.")
	return nil
}
