package core

import (
	"errors"
	"sync"
)

type Store interface {
	Get(string) (string, error)
	Put(string, string) error
	Delete(string) error
}

var ErrorNoSuchKey = errors.New("no such key")

type KeyValueStore struct {
	sync.RWMutex
	m map[string]string
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{m: make(map[string]string)}
}

func (k *KeyValueStore) Get(key string) (string, error) {
	k.RLock()
	value, ok := k.m[key]
	k.RUnlock()
	if ok {
		return value, nil
	}
	return "", ErrorNoSuchKey
}

func (k *KeyValueStore) Put(key string, value string) error {
	k.Lock()
	k.m[key] = value
	k.Unlock()
	return nil
}

func (k *KeyValueStore) Delete(key string) error {
	k.Lock()
	delete(k.m, key)
	k.Unlock()
	return nil
}
