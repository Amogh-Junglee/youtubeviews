package service

import (
	"context"
	"log"
	"time"
	"youtubeviews/models"
)

type MetricsService struct {
	service Service
}

func NewMetricsService(inner Service) *MetricsService {
	return &MetricsService{service: inner}
}

func (m *MetricsService) Increment(ctx context.Context, videoId string) (models.IncrementViewResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.Printf("Increment for VideoID %s took %v", videoId, duration)
	}()

	return m.service.Increment(ctx, videoId)
}

func (m *MetricsService) Get(ctx context.Context, videoId string) (models.ViewCountResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.Printf("Get for VideoID %s took %v", videoId, duration)
	}()

	return m.service.Get(ctx, videoId)
}

func (m *MetricsService) GetTopVideos(ctx context.Context, page int, limit int) (models.GetTopVideosResponse, error) {
	startTime := time.Now()

	resp, err := m.service.GetTopVideos(ctx, page, limit)

	duration := time.Since(startTime)
	logMetrics("GetTopVideos", duration, err)
	return resp, err
}

func logMetrics(methodName string, duration time.Duration, err error) {
	// TODO prometheus?
	log.Printf("Metrics: %s took %v, Error: %v", methodName, duration, err)
}
