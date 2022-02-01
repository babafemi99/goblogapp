package cache

import (
	"context"
	"encoding/json"
	"prodGo/entity"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) PostCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *redisCache) Set(key string, val *entity.Post) {
	client := cache.getClient()
	ctx := context.Background()

	json, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	client.Set(ctx, key, json, cache.expires*time.Second)
}
func (cache *redisCache) Get(key string) *entity.Post {
	client := cache.getClient()
	ctx := context.Background()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	post := entity.Post{}

	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		panic(err)
	}
	return &post
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}
