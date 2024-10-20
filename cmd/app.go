package main

import (
	"context"
	"fmt"
	"github.com/dafuqqqyunglean/todoRestAPI/config"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/middlewares"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/cache"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/auth"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/item"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/list"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

const (
	cacheKey = "todo_item:%d:%d"
	ttl      = time.Minute * 10
)

type App struct {
	ctx         utility.AppContext
	server      *api.Server
	cfg         config.Config
	repository  *sqlx.DB
	redisClient *redis.Client
}

func NewApp(ctx context.Context, logger *zap.SugaredLogger, cfg config.Config) *App {
	return &App{
		ctx: utility.NewAppContext(ctx, logger),
		cfg: cfg,
	}
}

func (a *App) Run() error {
	if err := a.server.Run(); err != nil {
		a.ctx.Logger.Fatalf("error occured while running http server: %s", err.Error())
	}

	a.ctx.Logger.Info("server is running")
	return nil
}

func (a *App) Shutdown() error {
	err := a.server.Shutdown(a.ctx.Ctx)
	if err != nil {
		a.ctx.Logger.Errorf("Failed to disconnect from server %v", err)
		return err
	}

	err = a.repository.Close()
	if err != nil {
		a.ctx.Logger.Errorf("failed to disconnect from bd %v", err)
	}

	a.ctx.Logger.Info("server shut down successfully")
	return nil
}

func (a *App) InitDatabase() error {
	var err error
	a.repository, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		a.cfg.Postgres.Host, a.cfg.Postgres.Port, a.cfg.Postgres.Username, a.cfg.Postgres.DBName, a.cfg.Postgres.Password, a.cfg.Postgres.SSLMode))
	if err != nil {
		a.ctx.Logger.Fatalf("failed to connect to postgres: %v", err)
	}

	err = a.repository.Ping()
	if err != nil {
		a.ctx.Logger.Fatalf("failed to ping database: %v", err)
	}

	a.redisClient = redis.NewClient(&redis.Options{
		Addr:     a.cfg.Redis.Address,
		Password: a.cfg.Redis.Password,
		DB:       a.cfg.Redis.DB,
	})

	_, err = a.redisClient.Ping(a.ctx.Ctx).Result()
	if err != nil {
		a.ctx.Logger.Fatalf("failed to ping Redis: %v", err)
	}

	return nil
}

func (a *App) InitService() {
	authService := auth.NewAuthorizationService(repository.NewAuthorizationPostgres(a.repository), a.ctx.Ctx)
	UserAuthMiddleware := middlewares.NewUserAuthMiddleware(authService)
	todoLists := list.NewTodoListService(repository.NewTodoListPostgres(a.repository), cache.NewRedisCache(a.redisClient, cacheKey, ttl))
	todoItems := item.NewTodoItemService(repository.NewTodoItemPostgres(a.repository), todoLists, cache.NewRedisCache(a.redisClient, cacheKey, ttl))
	a.server = api.NewServer(a.ctx, UserAuthMiddleware)
	a.server.HandleAuth(authService)
	a.server.HandleLists(a.ctx, todoLists)
	a.server.HandleItems(a.ctx, todoItems)
}
