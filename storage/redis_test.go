package storage

import (
	"os"
	"testing"
)

func TestNewRedisStorage_WithoutRedis(t *testing.T) {
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")

	_, err := NewRedisStorage()
	if err == nil {
		t.Error("Expected error when Redis is not available")
	}
}

func TestNewRedisStorage_InvalidPort(t *testing.T) {
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "invalid")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")

	_, err := NewRedisStorage()
	if err == nil {
		t.Error("Expected error with invalid port")
	}
}

func TestNewRedisStorage_InvalidDB(t *testing.T) {
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "invalid")

	_, err := NewRedisStorage()
	if err == nil {
		t.Error("Expected error with invalid DB")
	}
}
