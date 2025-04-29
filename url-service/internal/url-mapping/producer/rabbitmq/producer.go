package rabbitmq

import (
	"context"
	"encoding/json"
	"url-service/domain/events"
	"url-service/domain/interfaces"

	"github.com/rabbitmq/amqp091-go"
)

type urlMappingEventProducer struct {
	ch                   *amqp091.Channel
	exchange             string
	urlCreatedRoutingKey string
	urlUpdatedRoutingKey string
	urlDeletedRoutingKey string
}

func NewURLMappingEventProducer(ch *amqp091.Channel, exchange, urlCreatedRoutingKey, urlUpdatedRoutingKey, urlDeletedRoutingKey string) interfaces.URLMappingEventProducer {
	return &urlMappingEventProducer{
		ch:                   ch,
		exchange:             exchange,
		urlCreatedRoutingKey: urlCreatedRoutingKey,
		urlUpdatedRoutingKey: urlUpdatedRoutingKey,
		urlDeletedRoutingKey: urlDeletedRoutingKey,
	}
}

func (p *urlMappingEventProducer) ProduceURLMappingCreatedEvent(ctx context.Context, in events.URLCreatedEvent) error {
	payload, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = p.ch.PublishWithContext(ctx,
		p.exchange,             // exchange
		p.urlCreatedRoutingKey, // routing key
		false,                  // mandatory
		false,                  // immediate
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

func (p *urlMappingEventProducer) ProduceURLMappingDeletedEvent(ctx context.Context, in events.URLDeletedEvent) error {
	payload, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = p.ch.PublishWithContext(ctx,
		p.exchange,             // exchange
		p.urlDeletedRoutingKey, // routing key
		false,                  // mandatory
		false,                  // immediate
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

func (p *urlMappingEventProducer) ProduceURLMapingUpdatedEvent(ctx context.Context, in events.URLUpdatedEvent) error {
	payload, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = p.ch.PublishWithContext(ctx,
		p.exchange,             // exchange
		p.urlUpdatedRoutingKey, // routing key
		false,                  // mandatory
		false,                  // immediate
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
