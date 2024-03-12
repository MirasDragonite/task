package controllers

import (
	"context"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedis() (*cache.Cache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Password (if used)
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	myCache := cache.New(&cache.Options{
		Redis: client,
	})
	return myCache, nil
}
