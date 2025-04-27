package token_manager

import (
	"auth-service/domain/interfaces"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type revocationStore struct {
	client *redis.Client
}

func NewRevocationStore(client *redis.Client) interfaces.RevocationStore {
	return &revocationStore{client: client}
}

func (r *revocationStore) IsRevoked(jti string) bool {
	val, err := r.client.Exists(context.Background(), "blacklist:jti:"+jti).Result()
	if err != nil {
		return false
	}
	return val == 1
}

func (r *revocationStore) RevokeJTI(jti string, ttl time.Duration) error {
	return r.client.Set(context.Background(), fmt.Sprintf("blacklist:jti:%s", jti), "1", ttl).Err()
}

func (r *revocationStore) StoreRefreshToken(token, userUUID string, ttl time.Duration) error {
	return r.client.Set(context.Background(), fmt.Sprintf("refresh:%s", token), userUUID, ttl).Err()
}

func (r *revocationStore) DeleteRefreshToken(token string) error {
	return r.client.Del(context.Background(), fmt.Sprintf("refresh:%s", token)).Err()
}

func (r *revocationStore) GetRefreshToken(token string) (string, error) {
	val, err := r.client.Get(context.Background(), fmt.Sprintf("refresh:%s", token)).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("refresh token not found")
	}
	return val, err
}
