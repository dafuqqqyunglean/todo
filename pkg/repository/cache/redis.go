package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	client   *redis.Client
	cacheKey string
	ttl      time.Duration
}

func NewRedisCache(client *redis.Client, cacheKey string, ttl time.Duration) RedisCache {
	return RedisCache{
		client:   client,
		cacheKey: cacheKey,
		ttl:      ttl,
	}
}

func (r *RedisCache) SetItem(ctx utility.AppContext, userId, itemId int, item todo.TodoItem) {
	cacheKey := fmt.Sprintf(r.cacheKey, userId, itemId)

	itemJSON, _ := json.Marshal(item)

	r.client.Set(ctx.Ctx, cacheKey, itemJSON, r.ttl)
}

func (r *RedisCache) GetItem(ctx utility.AppContext, userId, itemId int) (todo.TodoItem, error) {
	cacheKey := fmt.Sprintf(r.cacheKey, userId, itemId)
	var item todo.TodoItem

	cachedItem, err := r.client.Get(ctx.Ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return item, fmt.Errorf("item not found in cache")
	} else if err != nil {
		return item, err
	}

	if err := json.Unmarshal([]byte(cachedItem), &item); err != nil {
		return item, err
	}

	return item, nil
}

func (r *RedisCache) SetList(ctx utility.AppContext, userId, itemId int, list todo.TodoList) {
	cacheKey := fmt.Sprintf(r.cacheKey, userId, itemId)

	itemJSON, _ := json.Marshal(list)

	r.client.Set(ctx.Ctx, cacheKey, itemJSON, r.ttl)
}

func (r *RedisCache) GetList(ctx utility.AppContext, userId, listId int) (todo.TodoList, error) {
	cacheKey := fmt.Sprintf(r.cacheKey, userId, listId)
	var list todo.TodoList

	cachedList, err := r.client.Get(ctx.Ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return list, fmt.Errorf("item not found in cache")
	} else if err != nil {
		return list, err
	}

	if err = json.Unmarshal([]byte(cachedList), &list); err == nil {
		return list, nil
	}

	return list, err
}

func (r *RedisCache) Delete(ctx utility.AppContext, userId, key int) {
	r.client.Del(ctx.Ctx, fmt.Sprintf(r.cacheKey, userId, key))
}
