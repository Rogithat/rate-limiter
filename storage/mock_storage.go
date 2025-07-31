package storage

import (
	"context"
	"sync"
)

type MockStorage struct {
	counters map[string]int64
	blocked  map[string]bool
	expiry   map[string]int
	mutex    sync.RWMutex
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		counters: make(map[string]int64),
		blocked:  make(map[string]bool),
		expiry:   make(map[string]int),
	}
}

func (m *MockStorage) Increment(ctx context.Context, key string) (int64, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.counters[key]++
	return m.counters[key], nil
}

func (m *MockStorage) SetExpiration(ctx context.Context, key string, duration int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.expiry[key] = duration
	return nil
}

func (m *MockStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if count, exists := m.counters[key]; exists {
		return count, nil
	}
	return 0, nil
}

func (m *MockStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.blocked[key], nil
}

func (m *MockStorage) SetBlocked(key string, blocked bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.blocked[key] = blocked
}

func (m *MockStorage) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.counters = make(map[string]int64)
	m.blocked = make(map[string]bool)
	m.expiry = make(map[string]int)
}
