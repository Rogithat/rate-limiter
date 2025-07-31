package storage

import (
	"context"
	"testing"
)

func TestNewMockStorage(t *testing.T) {
	mock := NewMockStorage()

	if mock == nil {
		t.Error("MockStorage should not be nil")
	}
	if mock.counters == nil {
		t.Error("counters map should be initialized")
	}
	if mock.blocked == nil {
		t.Error("blocked map should be initialized")
	}
	if mock.expiry == nil {
		t.Error("expiry map should be initialized")
	}
}

func TestMockStorage_Increment(t *testing.T) {
	mock := NewMockStorage()
	ctx := context.Background()
	key := "test-key"

	count, err := mock.Increment(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	count, err = mock.Increment(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	count, err = mock.Increment(ctx, "another-key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestMockStorage_SetExpiration(t *testing.T) {
	mock := NewMockStorage()
	ctx := context.Background()
	key := "test-key"
	duration := 300

	err := mock.SetExpiration(ctx, key, duration)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if mock.expiry[key] != duration {
		t.Errorf("Expected expiry %d, got %d", duration, mock.expiry[key])
	}
}

func TestMockStorage_GetCounter(t *testing.T) {
	mock := NewMockStorage()
	ctx := context.Background()
	key := "test-key"

	count, err := mock.GetCounter(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	mock.Increment(ctx, key)
	count, err = mock.GetCounter(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestMockStorage_IsBlocked(t *testing.T) {
	mock := NewMockStorage()
	ctx := context.Background()
	key := "test-key"

	blocked, err := mock.IsBlocked(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if blocked {
		t.Error("Expected not blocked")
	}

	mock.SetBlocked(key, true)
	blocked, err = mock.IsBlocked(ctx, key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !blocked {
		t.Error("Expected blocked")
	}
}

func TestMockStorage_SetBlocked(t *testing.T) {
	mock := NewMockStorage()
	key := "test-key"

	mock.SetBlocked(key, true)
	if !mock.blocked[key] {
		t.Error("Expected key to be blocked")
	}

	mock.SetBlocked(key, false)
	if mock.blocked[key] {
		t.Error("Expected key to not be blocked")
	}
}

func TestMockStorage_Reset(t *testing.T) {
	mock := NewMockStorage()
	ctx := context.Background()
	key := "test-key"

	mock.Increment(ctx, key)
	mock.SetExpiration(ctx, key, 300)
	mock.SetBlocked(key, true)

	count, _ := mock.GetCounter(ctx, key)
	if count != 1 {
		t.Error("Expected count to be 1 before reset")
	}

	mock.Reset()

	count, _ = mock.GetCounter(ctx, key)
	if count != 0 {
		t.Error("Expected count to be 0 after reset")
	}

	blocked, _ := mock.IsBlocked(ctx, key)
	if blocked {
		t.Error("Expected key to not be blocked after reset")
	}
}
