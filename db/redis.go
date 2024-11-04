package db

import (
	"YoutubeViews/models"
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) DbRepo {
	return &RedisRepo{client: client}
}

func (r *RedisRepo) Increment(ctx context.Context, req models.IncrementPayload) (models.IncrementViewResponse, error) {
	// Increment video views in Redis
	err := r.client.ZIncrBy(ctx, "video_views", 1, req.VideoID).Err()
	if err != nil {
		return models.IncrementViewResponse{}, err
	}

	// Fetch the updated score
	count, err := r.client.ZScore(ctx, "video_views", req.VideoID).Result()
	if err != nil {
		return models.IncrementViewResponse{}, err
	}

	return models.IncrementViewResponse{
		Views:   int(count),
		Increment: 1,
	}, nil
}

func (r *RedisRepo) Get(ctx context.Context, req models.ViewCountPayload) (models.ViewCountResponse, error) {
	// Fetch view count from Redis
	count, err := r.client.ZScore(ctx, "video_views", req.VideoID).Result()
	if err == redis.Nil {
		return models.ViewCountResponse{}, nil // video not found
	} else if err != nil {
		return models.ViewCountResponse{}, err
	}

	return models.ViewCountResponse{
		Views:   int(count),
	}, nil
}

func (r *RedisRepo) GetTopVideos(ctx context.Context, req models.GetTopVideosPayload) (models.GetTopVideosResponse, error) {
	// Fetch top videos from Redis
	start := (req.Page - 1) * req.Limit
	end := start + req.Limit - 1

	topVideos, err := r.client.ZRevRangeWithScores(ctx, "video_views", int64(start), int64(end)).Result()
	if err != nil {
		return models.GetTopVideosResponse{}, err
	}

	// Format response
	topVideoList := make([]map[string]interface{}, 0)
	for _, video := range topVideos {
		topVideoList = append(topVideoList, map[string]interface{}{
			"videoId": video.Member,
			"views":   video.Score,
		})
	}

	return models.GetTopVideosResponse{TopVideos: topVideoList}, nil
}
