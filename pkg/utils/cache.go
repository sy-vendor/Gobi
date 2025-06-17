package utils

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

var QueryCache *cache.Cache

func InitQueryCache(defaultExpiration, cleanupInterval time.Duration) {
	QueryCache = cache.New(defaultExpiration, cleanupInterval)
}

func GetQueryCache(key string) (interface{}, bool) {
	return QueryCache.Get(key)
}

func SetQueryCache(key string, value interface{}, ttl time.Duration) {
	QueryCache.Set(key, value, ttl)
}

func DeleteQueryCache(key string) {
	QueryCache.Delete(key)
}
