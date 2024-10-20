package main

import (
	"context"
	"fmt"
	"github.com/dafuqqqyunglean/todoRestAPI/config"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	prdLogger, _ := zap.NewProduction()
	defer prdLogger.Sync()
	logger := prdLogger.Sugar()

	fmt.Println(logger.Level())

	mainCtx := context.Background()
	ctx, cancel := context.WithCancel(mainCtx)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalf("failed to read settings: %s", err.Error())
	}

	app := NewApp(ctx, logger, cfg)
	if err := app.InitDatabase(); err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	app.InitService()

	if err = app.Run(); err != nil {
		logger.Errorf(err.Error())
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	if err = app.Shutdown(); err != nil {
		logger.Errorf(err.Error())
		return
	}
}
