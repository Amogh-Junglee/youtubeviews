package service

import (
	models "YoutubeViews/models"
	"context"
)

type Service interface {
	Increment(ctx context.Context, req models.IncrementPayload) (models.IncrementViewResponse, error)
	Get(ctx context.Context, req models.ViewCountPayload) (models.ViewCountResponse, error)
	GetTopVideos(ctx context.Context, req models.GetTopVideosPayload) (models.GetTopVideosResponse, error)
}
