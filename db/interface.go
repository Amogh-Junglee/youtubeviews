package db

import (
	"context"
	"youtubeviews/models"
)

type DbRepo interface {
	Increment(ctx context.Context, videoId string) (models.IncrementViewResponse, error)
	Get(ctx context.Context, videoId string) (models.ViewCountResponse, error)
	GetTopVideos(ctx context.Context, page int, limit int) (models.GetTopVideosResponse, error)
}
