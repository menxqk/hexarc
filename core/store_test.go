package core

import "testing"

func TestNew(t *testing.T) {
	store := NewKeyValueStore()
	if store == nil {
		t.Error("New store is nil")
	}

}
func TestGet(t *testing.T) {
	store := NewKeyValueStore()

	const key = "a-key"
	const value = "a-value"

	// Sanity check
	val, err := store.Get(key)
	if err == nil {
		t.Error("should return ErrorNoSuchKey for not inserted key")
	}
	if err != ErrorNoSuchKey {
		t.Errorf("error returned is not ErrorNoSuchKey: %q", err)
	}
	if val != "" {
		t.Errorf("should not return value %q for not inserted key", val)
	}

	err = store.Put(key, value)
	if err != nil {
		t.Error(err)
	}

	val, err = store.Get(key)
	if err != nil {
		t.Errorf("key was inserted, should not have returned error: %q", err)
	}
	if val != value {
		t.Errorf("val/value mismatch, val: %q, value: %q", val, value)
	}
}

func TestPut(t *testing.T) {
	store := NewKeyValueStore()

	const key1 = "a-key"
	const value1 = "a-value"
	const key2 = "b-key"
	const value2 = "b-value"

	err := store.Put(key1, value1)
	if err != nil {
		t.Error(err)
	}

	val, err := store.Get(key1)
	if err != nil {
		t.Error(err)
	}
	if val != value1 {
		t.Errorf("val/value mismatch, val: %q, value: %q", val, value1)
	}

	err = store.Put(key2, value2)
	if err != nil {
		t.Error(err)
	}

	val, err = store.Get(key2)
	if err != nil {
		t.Error(err)
	}
	if val != value2 {
		t.Errorf("val/value2 mismatch, val: %q, value: %q", val, value2)
	}
}

func TestDelete(t *testing.T) {
	store := NewKeyValueStore()

	const key1 = "a-key"
	const value1 = "a-value"
	const key2 = "b-key"
	const value2 = "b-value"
	const key3 = "c-key"
	const value3 = "c-value"

	err := store.Put(key1, value1)
	if err != nil {
		t.Error(err)
	}
	err = store.Put(key2, value2)
	if err != nil {
		t.Error(err)
	}
	err = store.Put(key3, value3)
	if err != nil {
		t.Error(err)
	}

	val, err := store.Get(key2)
	if err != nil {
		t.Error()
	}
	if val != value2 {
		t.Errorf("val/value2 mismatch, val: %q, value: %q", val, value2)
	}

	err = store.Delete(key2)
	if err != nil {
		t.Error(err)
	}
	val, err = store.Get(key2)
	if err == nil {
		t.Error("should have returned error for deleted key")
	}
	if err != ErrorNoSuchKey {
		t.Error("should have return ErrorNoSuchKey for deleted key")
	}
	if val != "" {
		t.Errorf("should not return value %q for deleted key", val)
	}

	err = store.Delete(key1)
	if err != nil {
		t.Error(err)
	}
	err = store.Delete(key3)
	if err != nil {
		t.Error(err)
	}
}
