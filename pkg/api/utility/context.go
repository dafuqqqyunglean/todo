package utility

import (
	"context"
	"go.uber.org/zap"
)

type AppContext struct {
	Ctx    context.Context
	Logger *zap.SugaredLogger
}

func NewAppContext(ctx context.Context, logger *zap.SugaredLogger) AppContext {
	return AppContext{
		Ctx:    ctx,
		Logger: logger,
	}
}
