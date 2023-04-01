package exchange

import (
	"context"
	"sync"
)

type inmemoryStore struct {
	data sync.Map
}

func newInmemoryStore() *inmemoryStore {
	return &inmemoryStore{}
}

func (s *inmemoryStore) Set(ctx context.Context, key string, value []byte) error {
	s.data.Store(key, value)
	return nil
}

func (s *inmemoryStore) Get(ctx context.Context, key string) ([]byte, error) {
	value, exists := s.data.Load(key)
	if !exists {
		return nil, ErrSecretNotFound
	}

	data, _ := value.([]byte)
	return data, nil
}
