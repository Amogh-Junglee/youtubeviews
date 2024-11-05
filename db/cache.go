package db

import (
	"context"
	"log"
)

type CacheRepo struct {
	cache    Cache[string, any]
	delegate DbRepo
}

func NewCacheRepo(delegate DbRepo, capacity int) DbRepo {
	return &CacheRepo{
		cache:    NewCache[string, any](capacity),
		delegate: delegate,
	}
}

func (c *CacheRepo) Increment(
	ctx context.Context,
	videoId string,
) (views int, increment int, err error) {
	// Get from cache
	if cachedResp, ok := c.cache.Get(videoId); ok {
		// Also increment in Redis
		redisViews, redisInc, err := c.delegate.Increment(ctx, videoId)
		if err != nil {
			return -1, -1, err
		}
		if viewCount, ok := cachedResp.(int); ok {
			// If cache hit, increment the count in the cache itself
			viewCount++
			c.cache.Put(videoId, viewCount)
			return viewCount, 1, nil
		}
		log.Printf("Cache hit but failed to type assert cached response to int")
		return redisViews, redisInc, nil
	}

	// If cache miss, delegate to the RedisRepo
	views, increment, err = c.delegate.Increment(ctx, videoId)
	if err != nil {
		return -1, -1, err
	}
	// Store the result in cache
	c.cache.Put(videoId, views)
	return views, increment, nil
}

func (c *CacheRepo) Get(ctx context.Context, videoId string) (views int, err error) {
	// Get from cache
	if cachedResp, ok := c.cache.Get(videoId); ok {
		if viewCount, ok := cachedResp.(int); ok {
			return viewCount, nil
		}
		// Handle if type assertion fails
		log.Printf("Cache hit but failed to type assert cached response to int")
	}

	// If cache miss, delegate to the RedisRepo
	redisViews, err := c.delegate.Get(ctx, videoId)
	if err != nil {
		return -1, err
	}

	// Store the result in cache
	c.cache.Put(videoId, redisViews)
	return redisViews, nil
}

// GetTopVideos retrieves the top videos based on the provided payload.
func (c *CacheRepo) GetTopVideos(ctx context.Context, page int, limit int) (topVideos []map[string]interface{}, err error) {
	return c.delegate.GetTopVideos(ctx, page, limit)
}
