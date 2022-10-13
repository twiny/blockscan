package inmemory

import "sync"

// InMemory -
type InMemory struct {
	mu *sync.Mutex
	db map[string]any
}
