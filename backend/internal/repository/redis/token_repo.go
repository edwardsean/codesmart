package redis

import (
	"context"
	"time"
)

type TokenRepository struct {
	client *RedisClient
}

func NewTokenRepository(client *RedisClient) *TokenRepository {
	return &TokenRepository{client: client}
}

func (r *TokenRepository) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return r.client.client.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}

func (r *TokenRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := r.client.client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil

}
