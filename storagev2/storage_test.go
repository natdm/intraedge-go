package storagev2_test

import (
	storage "coffeeserver/storagev2"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	store := storage.New()
	defer store.Close()
	err := store.Add("Nate", "Hyland")
	require.NoError(t, err)

	err = store.Add("", "Hyland")
	require.EqualError(t, err, "need a key")
}

func TestStorageImproved(t *testing.T) {
	tCases := []struct {
		name string
		err  error
		val  interface{}
		fn   func(*storage.Store) (interface{}, error)
	}{
		{
			name: "Add works fine",
			err:  nil,
			val:  nil,
			fn: func(s *storage.Store) (interface{}, error) {
				err := s.Add("Name", "Nate")
				return nil, err
			},
		},
		{
			name: "Add fails when called twice with the same key",
			err:  errors.New("key Name already exists"),
			val:  nil,
			fn: func(s *storage.Store) (interface{}, error) {
				err := s.Add("Name", "Nate")
				err = s.Add("Name", "Nate")
				return nil, err
			},
		},
		{
			name: "Add passes if the store vaidates the type against the custom validator",
			err:  nil,
			val:  nil,
			fn: func(s *storage.Store) (interface{}, error) {
				// set up the validator function to check if it's a string,
				s.ValidateFn = func(i interface{}) bool {
					_, ok := i.(string)
					return ok
				}

				// add a string to see if it passes the validator
				err := s.Add("Name", "Nate")
				return nil, err
			},
		},
		{
			name: "Add fails if the store doesn't vaidate the type against the custom validator",
			err:  storage.ErrInvalid,
			val:  nil,
			fn: func(s *storage.Store) (interface{}, error) {
				// set up the validator function to check if it's an integer,
				// or any other type but what you're giving it.
				s.ValidateFn = func(i interface{}) bool {
					_, ok := i.(int)
					return ok
				}

				// add something other than what the validator is checking against and return the error
				err := s.Add("Name", "Nate")
				return nil, err
			},
		},
	}

	for i, tCase := range tCases {
		t.Run(tCases[i].name, func(t *testing.T) {
			s := storage.New()
			defer s.Close()
			res, err := tCase.fn(s)
			if tCases[i].err != nil {
				require.Error(t, err)
				require.EqualError(t, tCases[i].err, err.Error())
				return
			}
			require.Equal(t, res, tCases[i].val)
		})
	}
}

func ExampleStore() {
	st := storage.New()
	defer st.Close()
	err := st.Add("Nate", "Hyland")
	if err != nil {
		panic(err)
	}
}

func ExampleStore2() {
	st := storage.New()
	defer st.Close()
	err := st.Add("Nate", "Hyland")
	if err != nil {
		panic(err)
	}
}
