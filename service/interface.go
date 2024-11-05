package service

import (
	"context"
	"youtubeviews/models"
)

type Service interface {
	Increment(ctx context.Context, videoId string) (models.IncrementViewResponse, error)
	Get(ctx context.Context, videoId string) (models.ViewCountResponse, error)
	GetTopVideos(ctx context.Context, page int, limit int) (models.GetTopVideosResponse, error)
}
