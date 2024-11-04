package db

import (
	"context"
	"fmt"

	"YoutubeViews/models"
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

func (c *CacheRepo) Increment(ctx context.Context, req models.IncrementPayload) (models.IncrementViewResponse, error) {
	// Get from cache
	if cachedResp, ok := c.cache.Get(req.VideoID); ok {
		// Also increment in Redis
		_, err := c.delegate.Increment(ctx, req)
		if err != nil {
			return models.IncrementViewResponse{}, err
		}

		// If cache hit, increment the count in the cache itself
		cachedIncrement := cachedResp.(models.IncrementViewResponse)
		cachedIncrement.Views++
		c.cache.Put(req.VideoID, cachedIncrement)

		return cachedIncrement, nil
	}

	// If cache miss, delegate to the RedisRepo
	resp, err := c.delegate.Increment(ctx, req)
	if err != nil {
		return models.IncrementViewResponse{}, err
	}

	// Store the result in cache
	c.cache.Put(req.VideoID, resp)
	return resp, nil
}

func (c *CacheRepo) Get(ctx context.Context, req models.ViewCountPayload) (models.ViewCountResponse, error) {
	// Get from cache
	if cachedResp, ok := c.cache.Get(req.VideoID); ok {
		if viewCount, ok := cachedResp.(models.IncrementViewResponse); ok {
			fmt.Println("ViewCount: ", viewCount)
			return models.ViewCountResponse{Views: viewCount.Views}, nil
		}
		fmt.Println("Cache : ", cachedResp)
		// Handle error if type assertion fails
		return models.ViewCountResponse{}, fmt.Errorf("failed to typecast cached response to ViewCountResponse")
	}

	// If cache miss, delegate to the RedisRepo
	resp, err := c.delegate.Get(ctx, req)
	if err != nil {
		return models.ViewCountResponse{}, err
	}

	// Store the result in cache
	c.cache.Put(req.VideoID, resp)
	return resp, nil
}

func (c *CacheRepo) GetTopVideos(ctx context.Context, req models.GetTopVideosPayload) (models.GetTopVideosResponse, error) {
	return c.delegate.GetTopVideos(ctx, req)
}
