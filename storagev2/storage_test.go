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
