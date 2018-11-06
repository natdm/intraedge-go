package storage

import (
	"errors"
	"fmt"
	"sync"
)

// Store stores a key of type string, and any data as the value
type Store struct {
	sync.RWMutex
	data map[string]interface{}
}

// New returns a new *Store
func New() *Store {
	return &Store{
		data: make(map[string]interface{}),
	}
}

// Add adds key to store, and returns an error if the key already exists
func (s *Store) Add(key string, val interface{}) error {
	if key == "" {
		return errors.New("need a key")
	}

	s.Lock()
	defer s.Unlock()
	s.data[key] = val
	return nil
}

// Val returns a value of interface that's associated with the key
// errors if there is no key, or the key supplied is blank
func (s *Store) Val(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("need a key")
	}

	s.RLock()
	defer s.RUnlock()
	val, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("value for key %s does not exist", key)
	}

	return val, nil
}

// State returns the current store state
func (s *Store) State() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.data
}
