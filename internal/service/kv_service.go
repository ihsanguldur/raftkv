package service

import (
	"errors"
	"strings"

	"github.com/ihsanguldur/raftkv/internal/kv"
)

var (
	ErrKeyRequired   = errors.New("key is required")
	ErrValueRequired = errors.New("value is required")
	ErrKeyNotFound   = errors.New("key not found")
)

type KVService struct {
	store *kv.Store
}

func NewKVService(store *kv.Store) *KVService {
	return &KVService{
		store: store,
	}
}

func (s *KVService) Get(key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", ErrKeyRequired
	}

	v, ok := s.store.Get(key)
	if !ok {
		return "", ErrKeyNotFound
	}

	return v, nil
}

func (s *KVService) Put(key string, value string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrKeyRequired
	}

	if value == "" {
		return ErrValueRequired
	}

	s.store.Set(key, value)
	return nil
}

func (s *KVService) Delete(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrKeyRequired
	}

	s.store.Delete(key)
	return nil
}
