package redis

import (
	"context"
	"errors"
	"log"
	"redirect-service/domain/entities"
	"redirect-service/domain/interfaces"
	"time"

	"github.com/redis/go-redis/v9"
)

type redirectURLCache struct {
	client *redis.Client
}

func NewRedirectURLCache(client *redis.Client) interfaces.RedirectURLCache {
	return &redirectURLCache{client: client}
}

func (c *redirectURLCache) Cache(ctx context.Context, shortCode, originalURL string) error {
	const expirationTime = 24 * time.Hour
	err := c.client.Set(ctx, shortCode, originalURL, expirationTime).Err()
	if err != nil {
		log.Println("Error caching redirect URL:", err)
		return entities.ErrInternalServer
	}
	return nil
}

func (c *redirectURLCache) Get(ctx context.Context, key string) (*string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, entities.ErrURLNotFound
		}
		log.Println("Error fetching from cache:", err)
		return nil, entities.ErrInternalServer
	}

	return &val, nil
}
