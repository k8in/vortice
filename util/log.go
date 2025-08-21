package util

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	mux                = sync.RWMutex{}
	logger *zap.Logger = defaultLogger()
	// 可替换的构建钩子，测试时注入失败场景以覆盖 fallback 分支
	buildLogger = func(cfg zap.Config) (*zap.Logger, error) {
		return cfg.Build()
	}
)

// defaultLogger returns a zap.Logger configured to log to standard output.
func defaultLogger() *zap.Logger {
	cfg := zap.Config{
		Encoding:          "console",
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		DisableStacktrace: true,
	}
	// 使用可替换构建函数，便于测试覆盖错误分支
	l, err := buildLogger(cfg)
	if err != nil {
		// fallback to zap.NewNop if config fails
		return zap.NewNop()
	}
	return l
}

// Logger returns a zap.Logger instance for logging.
// If logger is nil, returns defaultLogger to avoid panic.
func Logger() *zap.Logger {
	mux.RLock()
	defer mux.RUnlock()
	if logger == nil {
		return defaultLogger()
	}
	return logger
}

// SetLogger updates the global logger to the provided zap.Logger instance, ensuring thread safety.
func SetLogger(l *zap.Logger) {
	mux.Lock()
	defer mux.Unlock()
	logger = l
}
