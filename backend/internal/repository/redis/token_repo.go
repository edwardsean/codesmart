package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/domain"
)

type TokenRepository struct {
	client *RedisClient
}

type OAuthCodeData struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
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

func (r *TokenRepository) StoreOAuthCode(ctx context.Context, code string, data *OAuthCodeData) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.client.Set(ctx, "oauth_code:"+code, bytes, 60*time.Second).Err()
}

func (r *TokenRepository) ExchangeOAuthCode(ctx context.Context, code string) (*OAuthCodeData, error) {
	val, err := r.client.client.GetDel(ctx, "oauth_code:"+code).Result() // GetDel = get + delete atomically (one-time use)
	if err != nil {
		return nil, err
	}

	var data OAuthCodeData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}

	return &data, nil
}
