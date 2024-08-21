package profiles

import (
	"bytes"
	"encoding/json"
)

type MemoryStore struct {
	namespace string
	key       string

	memory *map[string]interface{}
}

// NewMemoryStore creates a new in-memory store
// JSON is used to serialize the data to ensure the interface is consistent with other store implementations
var NewMemoryStore NewStoreInterface = func(namespace string, key string) StoreInterface {
	memory := make(map[string]interface{})
	return &MemoryStore{
		namespace: namespace,
		key:       key,
		memory:    &memory,
	}
}

func (k *MemoryStore) Exists() bool {
	m := *k.memory
	_, ok := m[k.key]
	return ok
}

func (k *MemoryStore) Get(value interface{}) error {
	m := *k.memory
	v, ok := m[k.key]
	if !ok {
		return nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewReader(b)).Decode(value)
}

func (k *MemoryStore) Set(value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m := *k.memory
	m[k.key] = b
	// maybe write back to k.memory
	// k.memory = &m
	return nil
}

func (k *MemoryStore) Delete() error {
	m := *k.memory
	delete(m, k.key)
	// maybe write back to k.memory
	return nil
}
