package object

import (
	"sync/atomic"
)

type Registry struct {
	entries   map[string][]*Definition
	factories map[string]string
	readonly  *atomic.Bool
}

func newRegistry() *Registry {
	readonly := &atomic.Bool{}
	readonly.Store(false)
	return &Registry{
		entries:   map[string][]*Definition{},
		factories: map[string]string{},
		readonly:  readonly,
	}
}
