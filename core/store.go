package core

import (
	"errors"
	"log"
	"sync"

	"github.com/menxqk/hexarc/backend"
)

var ErrorNoSuchKey = errors.New("no such key")

type KeyValueStore struct {
	sync.RWMutex
	m        map[string]string
	transact backend.TransactionLogger
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		m:        make(map[string]string),
		transact: zeroTransactionLogger{},
	}
}

func (k *KeyValueStore) WithTransactionLogger(tl backend.TransactionLogger) *KeyValueStore {
	k.transact = tl
	return k
}

func (k *KeyValueStore) Init() error {
	var err error

	events, errors := k.transact.ReadEvents()
	count, ok, e := 0, true, backend.Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errors:

		case e, ok = <-events:

			switch e.EventType {
			case backend.EventPut:
				k.Lock()
				k.m[e.Key] = e.Value
				k.Unlock()
				count++
			case backend.EventDelete:
				k.Lock()
				delete(k.m, e.Key)
				k.Unlock()
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)

	k.transact.Run()

	return err
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
	k.transact.WritePut(key, value)
	return nil
}

func (k *KeyValueStore) Delete(key string) error {
	k.Lock()
	delete(k.m, key)
	k.Unlock()
	k.transact.WriteDelete(key)
	return nil
}

func (k *KeyValueStore) Size() uint64 {
	k.Lock()
	size := uint64(len(k.m))
	k.Unlock()
	return size
}

func (k *KeyValueStore) WaitForTransactionLogger() {
	k.transact.Wait()
}

type zeroTransactionLogger struct{}

func (z zeroTransactionLogger) WritePut(key, value string)                       {}
func (z zeroTransactionLogger) WriteDelete(key string)                           {}
func (z zeroTransactionLogger) Err() <-chan error                                { return nil }
func (z zeroTransactionLogger) LastSequence() uint64                             { return 0 }
func (z zeroTransactionLogger) Run()                                             {}
func (z zeroTransactionLogger) Wait()                                            {}
func (z zeroTransactionLogger) Close() error                                     { return nil }
func (z zeroTransactionLogger) ReadEvents() (<-chan backend.Event, <-chan error) { return nil, nil }
