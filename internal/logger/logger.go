package logger

import "go.uber.org/zap"

func NewLogger(isDev bool) (*zap.Logger, error) {
	if (isDev) {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}