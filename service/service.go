package service

import (
	"context"
	"fmt"

	"youtubeviews/db"
	"youtubeviews/models"

	"github.com/redis/go-redis/v9"
)

type VideoService struct {
	repo db.DbRepo
}

func CreateNewVideoService(redisClient *redis.Client, capacity int) *VideoService {
	// Create Redis
	redisRepo := db.NewRedisRepo(redisClient)

	// Create Cache
	cacheRepo := db.NewCacheRepo(redisRepo, capacity)

	return &VideoService{repo: cacheRepo}
}

func (videoService *VideoService) Increment(ctx context.Context, videoId string) (models.IncrementViewResponse, error) {
	// Validate that the VideoID is provided
	if videoId == "" {
		return models.IncrementViewResponse{}, fmt.Errorf("videoId is required")
	}

	return videoService.repo.Increment(ctx, videoId)
}

func (videoService *VideoService) Get(ctx context.Context, videoId string) (models.ViewCountResponse, error) {
	if videoId == "" {
		return models.ViewCountResponse{}, fmt.Errorf("videoId is required")
	}

	return videoService.repo.Get(ctx, videoId)
}

func (videoService *VideoService) GetTopVideos(ctx context.Context, page int, limit int) (models.GetTopVideosResponse, error) {
	if page < 1{
		return models.GetTopVideosResponse{}, fmt.Errorf("page must be greater than 0")
	}
	return videoService.repo.GetTopVideos(ctx, page, limit)
}
