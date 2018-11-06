package storage_test

import (
	"coffeeserver/storage"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	store := storage.New()
	err := store.Add("Nate", "Hyland")
	require.NoError(t, err)

	err = store.Add("", "Hyland")
	require.EqualError(t, err, "need a key")
}

func TestStorageConcurrency(t *testing.T) {
	store := storage.New()

	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			err := store.Add(string(i), i)
			require.NoError(t, err)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
