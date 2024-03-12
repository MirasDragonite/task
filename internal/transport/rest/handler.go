package rest

import (
	"miras/internal/services"

	"github.com/go-redis/cache/v9"
)

type Handler struct {
	Service *services.Service
	cache   *cache.Cache
}

func NewHandler(service *services.Service, cache *cache.Cache) *Handler {
	return &Handler{Service: service, cache: cache}
}
