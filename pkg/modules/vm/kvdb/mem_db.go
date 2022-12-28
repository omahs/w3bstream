package kvdb

import (
	"errors"
	"fmt"
)

type memDB struct {
	db map[string][]byte
}

func NewMemDB() *memDB {
	return &memDB{db: make(map[string][]byte)}
}

func (m *memDB) Get(key string) ([]byte, error) {
	value, ok := m.db[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("key[%s] not found", key))
	}
	return value, nil
}

func (m *memDB) Set(key string, value []byte) error {
	m.db[key] = value
	return nil
}
