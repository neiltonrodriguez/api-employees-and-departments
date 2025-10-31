package cache

import (
	"context"
	"sync"
	"time"
)

// MockCache is a test implementation of Cache interface for unit testing
// Usage example in tests:
//
//	mock := cache.NewMockCache()
//	service := department.NewService(repo, empRepo, logger, mock, 5*time.Minute)
//	// ... test service methods
//	// ... assert mock.Store contains expected entries
type MockCache struct {
	Store   map[string]string
	mu      sync.RWMutex
	GetCalls    int
	SetCalls    int
	DeleteCalls int
}

// NewMockCache creates a new mock cache for testing
func NewMockCache() *MockCache {
	return &MockCache{
		Store: make(map[string]string),
	}
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.GetCalls++
	return m.Store[key], nil
}

func (m *MockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SetCalls++
	m.Store[key] = value
	return nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DeleteCalls++
	delete(m.Store, key)
	return nil
}

func (m *MockCache) DeleteByPattern(ctx context.Context, pattern string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Simple implementation: clear all for testing
	m.Store = make(map[string]string)
	return nil
}

func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.Store[key]
	return exists, nil
}

func (m *MockCache) Ping(ctx context.Context) error {
	return nil
}

// Reset clears the mock cache (useful between test cases)
func (m *MockCache) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Store = make(map[string]string)
	m.GetCalls = 0
	m.SetCalls = 0
	m.DeleteCalls = 0
}

// HasKey checks if a key exists in cache
func (m *MockCache) HasKey(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.Store[key]
	return exists
}
