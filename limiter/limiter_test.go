package limiter

import (
	"context"
	"os"
	"testing"

	"rate-limiter/storage"
)

func TestNewConfig(t *testing.T) {
	os.Setenv("DEFAULT_IP_LIMIT", "5")
	os.Setenv("DEFAULT_IP_BLOCK_DURATION", "300")
	os.Setenv("DEFAULT_TOKEN_LIMIT", "10")
	os.Setenv("DEFAULT_TOKEN_BLOCK_DURATION", "300")

	config := NewConfig()

	if config.IPLimit == 0 {
		t.Error("IPLimit should not be zero")
	}
	if config.IPBlockDuration == 0 {
		t.Error("IPBlockDuration should not be zero")
	}
	if config.TokenLimit == 0 {
		t.Error("TokenLimit should not be zero")
	}
	if config.TokenBlockDuration == 0 {
		t.Error("TokenBlockDuration should not be zero")
	}

	if config.IPLimit != 5 {
		t.Errorf("Expected IPLimit 5, got %d", config.IPLimit)
	}
	if config.IPBlockDuration != 300 {
		t.Errorf("Expected IPBlockDuration 300, got %d", config.IPBlockDuration)
	}
	if config.TokenLimit != 10 {
		t.Errorf("Expected TokenLimit 10, got %d", config.TokenLimit)
	}
	if config.TokenBlockDuration != 300 {
		t.Errorf("Expected TokenBlockDuration 300, got %d", config.TokenBlockDuration)
	}
}

func TestNewRateLimiter(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	limiter := NewRateLimiter(mockStorage, config)

	if limiter == nil {
		t.Error("RateLimiter should not be nil")
	}
	if limiter.storage != mockStorage {
		t.Error("Storage should be set correctly")
	}
	if limiter.config != config {
		t.Error("Config should be set correctly")
	}
}

func TestCheckRateLimit_IPBased(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &Config{
		IPLimit:            3,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	limiter := NewRateLimiter(mockStorage, config)
	ctx := context.Background()
	ip := "192.168.1.1"

	limited, err := limiter.CheckRateLimit(ctx, ip, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("First request should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, ip, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("Second request should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, ip, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("Third request should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, ip, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !limited {
		t.Error("Fourth request should be limited")
	}
}

func TestCheckRateLimit_TokenBased(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         2,
		TokenBlockDuration: 300,
	}

	limiter := NewRateLimiter(mockStorage, config)
	ctx := context.Background()
	ip := "192.168.1.1"
	token := "test-token"

	limited, err := limiter.CheckRateLimit(ctx, ip, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("First request with token should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, ip, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("Second request with token should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, ip, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !limited {
		t.Error("Third request with token should be limited")
	}
}

func TestCheckRateLimit_IPBlocked(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	limiter := NewRateLimiter(mockStorage, config)
	ctx := context.Background()
	ip := "192.168.1.1"

	mockStorage.SetBlocked("ip:192.168.1.1", true)

	limited, err := limiter.CheckRateLimit(ctx, ip, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !limited {
		t.Error("Request should be limited when IP is blocked")
	}
}

func TestCheckRateLimit_TokenBlocked(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	limiter := NewRateLimiter(mockStorage, config)
	ctx := context.Background()
	ip := "192.168.1.1"
	token := "test-token"

	mockStorage.SetBlocked("token:test-token", true)

	limited, err := limiter.CheckRateLimit(ctx, ip, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !limited {
		t.Error("Request should be limited when token is blocked")
	}
}

func TestCheckRateLimit_TokenOverridesIP(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &Config{
		IPLimit:            2,
		IPBlockDuration:    300,
		TokenLimit:         5,
		TokenBlockDuration: 300,
	}

	limiter := NewRateLimiter(mockStorage, config)
	ctx := context.Background()
	ip := "192.168.1.1"
	token := "test-token"

	limited, err := limiter.CheckRateLimit(ctx, ip, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("First IP request should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, ip, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("Second IP request should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, ip, token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("First token request should not be limited")
	}
}

func TestCheckRateLimit_StorageError(t *testing.T) {
	mockStorage := &errorStorage{}
	config := &Config{
		IPLimit:            5,
		IPBlockDuration:    300,
		TokenLimit:         10,
		TokenBlockDuration: 300,
	}

	limiter := NewRateLimiter(mockStorage, config)
	ctx := context.Background()

	_, err := limiter.CheckRateLimit(ctx, "192.168.1.1", "")
	if err == nil {
		t.Error("Expected error when storage fails")
	}
}

type errorStorage struct{}

func (e *errorStorage) Increment(ctx context.Context, key string) (int64, error) {
	return 0, context.DeadlineExceeded
}

func (e *errorStorage) SetExpiration(ctx context.Context, key string, duration int) error {
	return context.DeadlineExceeded
}

func (e *errorStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	return 0, context.DeadlineExceeded
}

func (e *errorStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	return false, context.DeadlineExceeded
}

func TestCheckRateLimit_EdgeCases(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	config := &Config{
		IPLimit:            1,
		IPBlockDuration:    300,
		TokenLimit:         1,
		TokenBlockDuration: 300,
	}

	limiter := NewRateLimiter(mockStorage, config)
	ctx := context.Background()

	limited, err := limiter.CheckRateLimit(ctx, "", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("Empty IP should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if limited {
		t.Error("First request should not be limited")
	}

	limited, err = limiter.CheckRateLimit(ctx, "192.168.1.1", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !limited {
		t.Error("Second request should be limited")
	}
}
