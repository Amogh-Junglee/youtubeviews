package service

import (
	"context"

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

func (videoService *VideoService) Increment(ctx context.Context, req models.IncrementPayload) (models.IncrementViewResponse, error) {
	return videoService.repo.Increment(ctx, req)
}

func (videoService *VideoService) Get(ctx context.Context, req models.ViewCountPayload) (models.ViewCountResponse, error) {
	return videoService.repo.Get(ctx, req)
}

func (videoService *VideoService) GetTopVideos(ctx context.Context, req models.GetTopVideosPayload) (models.GetTopVideosResponse, error) {
	return videoService.repo.GetTopVideos(ctx, req)
}
