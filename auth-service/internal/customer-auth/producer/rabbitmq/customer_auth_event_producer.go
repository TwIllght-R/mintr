package rabbitmq

import (
	"auth-service/domain/events"
	"auth-service/domain/interfaces"
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type customerAuthEventProducer struct {
	ch *amqp091.Channel
}

func NewCustomerAuthEventProducer(ch *amqp091.Channel) interfaces.CustomerAuthEventProducer {
	return &customerAuthEventProducer{ch: ch}
}

func (p *customerAuthEventProducer) ProduceEmailVerificationEvent(ctx context.Context, in events.EmailVerificationEvent) error {
	payload, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = p.ch.PublishWithContext(ctx,
		"notification",             // exchange
		"email.verification.queue", // routing key
		false,                      // mandatory
		false,                      // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
