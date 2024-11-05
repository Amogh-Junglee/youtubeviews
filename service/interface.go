package service

import (
	"context"
)

type Service interface {
	Increment(ctx context.Context, videoId string) (views int, increment int, err error)
	Get(ctx context.Context, videoId string) (views int, err error)
	GetTopVideos(ctx context.Context, page int, limit int) (topVideos []map[string]interface{}, err error)
}
