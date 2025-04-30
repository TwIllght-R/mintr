package redis

import (
	"context"
	"errors"
	"fmt"
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

func (c *redirectURLCache) StoreMappingURL(ctx context.Context, shortCode, originalURLWithUTM string) error {
	const expirationTime = 24 * time.Hour
	err := c.client.Set(ctx, shortCode, originalURLWithUTM, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("failed to cache redirect URL: %w", err)
	}
	return nil
}

func (c *redirectURLCache) GetMappingURL(ctx context.Context, key string) (*string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, entities.ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to get value from cache: %w", err)
	}

	return &val, nil
}
