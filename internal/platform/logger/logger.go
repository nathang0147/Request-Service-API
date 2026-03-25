package logger

import "go.uber.org/zap"

func New(level string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	if err := config.Level.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	return config.Build()
}
