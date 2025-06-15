package engines

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_BasicCRUD(t *testing.T) {
	storage := NewMemoryStorage[string, int]()

	// Test Set
	storage.Mu.Lock()
	storage.Data["key1"] = 10
	storage.Mu.Unlock()

	// Test Get
	storage.Mu.RLock()
	val, ok := storage.Data["key1"]
	storage.Mu.RUnlock()

	assert.True(t, ok, "value should exist")
	assert.Equal(t, 10, val)

	// Test Update
	storage.Mu.Lock()
	storage.Data["key1"] = 20
	storage.Mu.Unlock()

	storage.Mu.RLock()
	val, ok = storage.Data["key1"]
	storage.Mu.RUnlock()

	assert.True(t, ok)
	assert.Equal(t, 20, val)

	// Test Delete
	storage.Mu.Lock()
	delete(storage.Data, "key1")
	storage.Mu.Unlock()

	storage.Mu.RLock()
	_, ok = storage.Data["key1"]
	storage.Mu.RUnlock()

	assert.False(t, ok, "value should have been deleted")
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	storage := NewMemoryStorage[string, int]()
	var wg sync.WaitGroup

	numGoroutines := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + string(rune(i))
			storage.Mu.Lock()
			storage.Data[key] = i
			storage.Mu.Unlock()
		}(i)
	}
	wg.Wait()

	// Validate all writes
	storage.Mu.RLock()
	assert.Len(t, storage.Data, numGoroutines)
	storage.Mu.RUnlock()
}
