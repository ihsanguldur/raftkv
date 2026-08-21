package service

import (
	"errors"
	"strings"

	"github.com/ihsanguldur/raftkv/internal/kv"
	"github.com/ihsanguldur/raftkv/internal/raft"
)

var (
	ErrKeyRequired   = errors.New("key is required")
	ErrValueRequired = errors.New("value is required")
	ErrKeyNotFound   = errors.New("key not found")
	ErrNotLeader     = errors.New("not leader")
	ErrTimeout       = errors.New("commit timeout")
)

type KVService struct {
	store *kv.Store
	node  *raft.Node
}

func NewKVService(store *kv.Store, node *raft.Node) *KVService {
	return &KVService{
		store: store,
		node:  node,
	}
}

func (s *KVService) LeaderAddr() string {
	return s.node.LeaderAddr()
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

	if err := s.node.Propose(raft.Command{Op: "put", Key: key, Value: value}); err != nil {
		if errors.Is(err, raft.ErrNotLeader) {
			return ErrNotLeader
		}
		if errors.Is(err, raft.ErrCommitTimeout) {
			return ErrTimeout
		}
		return err
	}
	return nil
}

func (s *KVService) Delete(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrKeyRequired
	}

	if err := s.node.Propose(raft.Command{Op: "delete", Key: key}); err != nil {
		if errors.Is(err, raft.ErrNotLeader) {
			return ErrNotLeader
		}
		if errors.Is(err, raft.ErrCommitTimeout) {
			return ErrTimeout
		}
		return err
	}
	return nil
}
