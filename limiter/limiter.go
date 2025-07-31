package limiter

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"rate-limiter/storage"
)

type Config struct {
	IPLimit            int
	IPBlockDuration    int
	TokenLimit         int
	TokenBlockDuration int
}

func NewConfig() *Config {
	ipLimit, _ := strconv.Atoi(os.Getenv("DEFAULT_IP_LIMIT"))
	ipBlockDuration, _ := strconv.Atoi(os.Getenv("DEFAULT_IP_BLOCK_DURATION"))
	tokenLimit, _ := strconv.Atoi(os.Getenv("DEFAULT_TOKEN_LIMIT"))
	tokenBlockDuration, _ := strconv.Atoi(os.Getenv("DEFAULT_TOKEN_BLOCK_DURATION"))

	return &Config{
		IPLimit:            ipLimit,
		IPBlockDuration:    ipBlockDuration,
		TokenLimit:         tokenLimit,
		TokenBlockDuration: tokenBlockDuration,
	}
}

type RateLimiter struct {
	storage storage.Storage
	config  *Config
}

func NewRateLimiter(storage storage.Storage, config *Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

func (rl *RateLimiter) CheckRateLimit(ctx context.Context, ip string, token string) (bool, error) {
	ipBlocked, err := rl.storage.IsBlocked(ctx, fmt.Sprintf("ip:%s", ip))
	if err != nil {
		return false, err
	}
	if ipBlocked {
		return true, nil
	}

	if token != "" {
		tokenBlocked, err := rl.storage.IsBlocked(ctx, fmt.Sprintf("token:%s", token))
		if err != nil {
			return false, err
		}
		if tokenBlocked {
			return true, nil
		}

		count, err := rl.storage.Increment(ctx, fmt.Sprintf("token:%s", token))
		if err != nil {
			return false, err
		}

		if count == 1 {
			err = rl.storage.SetExpiration(ctx, fmt.Sprintf("token:%s", token), rl.config.TokenBlockDuration)
			if err != nil {
				return false, err
			}
		}

		if count > int64(rl.config.TokenLimit) {
			err = rl.storage.SetExpiration(ctx, fmt.Sprintf("token:%s:blocked", token), rl.config.TokenBlockDuration)
			if err != nil {
				return false, err
			}
			return true, nil
		}

		return false, nil
	}

	count, err := rl.storage.Increment(ctx, fmt.Sprintf("ip:%s", ip))
	if err != nil {
		return false, err
	}

	if count == 1 {
		err = rl.storage.SetExpiration(ctx, fmt.Sprintf("ip:%s", ip), rl.config.IPBlockDuration)
		if err != nil {
			return false, err
		}
	}

	if count > int64(rl.config.IPLimit) {
		err = rl.storage.SetExpiration(ctx, fmt.Sprintf("ip:%s:blocked", ip), rl.config.IPBlockDuration)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}
