package db

import (
	"context"
	"log"

	"youtubeviews/models"
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

func (c *CacheRepo) Increment(ctx context.Context, videoId string) (models.IncrementViewResponse, error) {
	// Get from cache
	if cachedResp, ok := c.cache.Get(videoId); ok {
		// Also increment in Redis
		redisResponse, err := c.delegate.Increment(ctx, videoId)
		if err != nil {
			return models.IncrementViewResponse{}, err
		}
		if viewCount, ok := cachedResp.(int); ok{
			// If cache hit, increment the count in the cache itself
			cachedIncrement := models.IncrementViewResponse{Views: viewCount, Increment: 1}
			cachedIncrement.Views++
			c.cache.Put(videoId, cachedIncrement.Views)
			return cachedIncrement, nil
		}
		log.Printf("Cache hit but failed to type assert cached response to int")
		return redisResponse, nil
	}

	// If cache miss, delegate to the RedisRepo
	resp, err := c.delegate.Increment(ctx, videoId)
	if err != nil {
		return models.IncrementViewResponse{}, err
	}
	// Store the result in cache
	c.cache.Put(videoId, resp.Views)
	return resp, nil
}

func (c *CacheRepo) Get(ctx context.Context, videoId string) (models.ViewCountResponse, error) {
	// Get from cache
	if cachedResp, ok := c.cache.Get(videoId); ok {
		if viewCount, ok := cachedResp.(int); ok {
			return models.ViewCountResponse{Views: viewCount}, nil
		}
		// Handle if type assertion fails
		log.Printf("Cache hit but failed to type assert cached response to int")
	}

	// If cache miss, delegate to the RedisRepo
	resp, err := c.delegate.Get(ctx, videoId)
	if err != nil {
		return models.ViewCountResponse{}, err
	}

	// Store the result in cache
	c.cache.Put(videoId, resp.GetViews())
	return resp, nil
}

// GetTopVideos retrieves the top videos based on the provided payload.
func (c *CacheRepo) GetTopVideos(ctx context.Context, page int, limit int) (models.GetTopVideosResponse, error) {
	return c.delegate.GetTopVideos(ctx, page, limit)
}
