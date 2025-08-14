package util

import (
	"sync"

	"go.uber.org/zap"
)

var (
	mux    = sync.RWMutex{}
	logger = zap.L()
)

// Logger returns a zap.Logger instance for logging.
func Logger() *zap.Logger {
	mux.RLock()
	defer mux.RUnlock()
	return logger
}

// SetLogger updates the global logger to the provided zap.Logger instance, ensuring thread safety.
func SetLogger(l *zap.Logger) {
	mux.Lock()
	defer mux.Unlock()
	logger = l
}
