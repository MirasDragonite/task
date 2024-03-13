package controllers

import (
	"context"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedis() (*cache.Cache, error) {
	// connecting to our redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // Password (if used)
		DB:       0,
	})

	//check if there is any error
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	//creating new cache for storing some data
	myCache := cache.New(&cache.Options{
		Redis: client,
	})
	return myCache, nil
}
