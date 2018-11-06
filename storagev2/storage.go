package storagev2

import (
	"errors"
	"fmt"
)

// Store stores a key of type string, and any data as the value
type Store struct {
	data map[string]interface{}
	ops  chan func(map[string]interface{})
}

// New returns a new *Store, and the caller is responsible for calling Close()
func New() *Store {
	data := make(map[string]interface{})
	ops := make(chan func(map[string]interface{}))
	store := &Store{
		data: data,
		ops:  ops,
	}

	go func() {
		for op := range ops {
			op(data)
		}
	}()
	return store
}

func (s *Store) Close() error {
	close(s.ops)
	return nil
}

// Add adds key to store, and returns an error if the key already exists
func (s *Store) Add(key string, val interface{}) error {
	err := make(chan error)
	go func() {
		s.ops <- func(state map[string]interface{}) {
			if key == "" {
				err <- errors.New("need a key")
				return
			}
			if _, ok := state[key]; ok {
				err <- fmt.Errorf("key %s already exists", key)
				return
			}
			state[key] = val
			err <- nil
		}
	}()
	return <-err
}

// Val returns a value of interface that's associated with the key
// errors if there is no key, or the key supplied is blank
func (s *Store) Val(key string) (interface{}, error) {
	val, err := make(chan interface{}), make(chan error)
	go func() {
		s.ops <- func(state map[string]interface{}) {
			if key == "" {
				val <- nil
				err <- errors.New("need a key")
				return
			}

			v, ok := state[key]
			if !ok {
				val <- nil
				err <- fmt.Errorf("value for key %s does not exist", key)
				return
			}

			val <- v
			err <- nil
		}
	}()
	return <-val, <-err
}

// State returns the current store state
func (s *Store) State() map[string]interface{} {
	data := make(chan map[string]interface{})
	go func() {
		s.ops <- func(state map[string]interface{}) {
			data <- state
		}
	}()
	return <-data
}
