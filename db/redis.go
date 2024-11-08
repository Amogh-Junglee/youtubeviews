package db

import (
	"context"

	"math"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) DbRepo {
	return &RedisRepo{client: client}
}

func (r *RedisRepo) Increment(ctx context.Context, videoId string) (views int, increment int, err error) {
	// Increment video views in Redis
	err = r.client.ZIncrBy(ctx, "video_views", 1, videoId).Err()
	if err != nil {
		return -1, -1, err
	}

	// Fetch the updated score
	count, err := r.client.ZScore(ctx, "video_views", videoId).Result()
	if err != nil {
		return -1, -1, err
	}

	return int(math.Round(count)), 1, nil
}

func (r *RedisRepo) Get(ctx context.Context, videoId string) (views int, err error) {
	// Fetch view count from Redis
	count, err := r.client.ZScore(ctx, "video_views", videoId).Result()
	if err == redis.Nil {
		return 0, nil // video not found
	} else if err != nil {
		return -1, err
	}

	return int(count), nil
}

func (r *RedisRepo) GetTopVideos(ctx context.Context, page int, limit int) (topVideos []map[string]interface{}, err error) {
	// Fetch top videos from Redis
	start := (page - 1) * limit
	end := start + limit - 1

	redisTopVideos, err := r.client.ZRevRangeWithScores(ctx, "video_views", int64(start), int64(end)).Result()
	if err != nil {
		return make([]map[string]interface{}, 0), err
	}

	// Format response
	topVideoList := make([]map[string]interface{}, 0)
	for _, video := range redisTopVideos {
		topVideoList = append(topVideoList, map[string]interface{}{
			"videoId": video.Member,
			"views":   video.Score,
		})
	}

	return topVideoList, nil
}
