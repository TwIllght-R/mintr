package redis

import (
	"auth-service/domain/interfaces"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type customerAuthCache struct {
	client *redis.Client
}

func NewCustomerAuthCache(client *redis.Client) interfaces.CustomerAuthCache {
	return &customerAuthCache{client: client}
}

func (c *customerAuthCache) StoreEmailVerificationToken(ctx context.Context, email, token string) error {
	return c.client.Set(ctx, "email_verification:"+token, email, 15*time.Minute).Err()
}

func (c *customerAuthCache) VerifyEmailVerificationToken(ctx context.Context, token string) (*string, error) {
	val, err := c.client.Get(ctx, "email_verification:"+token).Result()
	if err != nil {
		return nil, err
	}
	return &val, nil

}
