package infrastructure

import (
	"time"

	"gobi/pkg/utils"
)

// CacheServiceImpl implements CacheService
type CacheServiceImpl struct{}

// NewCacheService creates a new CacheService instance
func NewCacheService() *CacheServiceImpl {
	return &CacheServiceImpl{}
}

// Get retrieves a value from cache
func (s *CacheServiceImpl) Get(key string) (interface{}, bool) {
	return utils.QueryCache.Get(key)
}

// Set sets a value in cache with smart TTL
func (s *CacheServiceImpl) Set(key string, value interface{}, ttl time.Duration) {
	utils.QueryCache.Set(key, value, ttl)
}

// Delete deletes a value from cache
func (s *CacheServiceImpl) Delete(key string) {
	utils.QueryCache.Delete(key)
}

// Flush flushes all cache
func (s *CacheServiceImpl) Flush() {
	utils.QueryCache.Flush()
}

// GetStats returns cache statistics
func (s *CacheServiceImpl) GetStats() map[string]interface{} {
	return utils.GetCacheStats()
}
