package engines

import (
	"sync"
)

type MemoryStorage[K comparable, V any] struct {
	Data map[K]V
	Mu   *sync.RWMutex
}

func NewMemoryStorage[K comparable, V any]() *MemoryStorage[K, V] {
	return &MemoryStorage[K, V]{
		Data: make(map[K]V),
		Mu:   &sync.RWMutex{},
	}
}
