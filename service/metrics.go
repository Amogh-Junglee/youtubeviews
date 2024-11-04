package service

import (
	"YoutubeViews/models"
	"context"
	"log"
	"time"
)

type MetricsService struct {
	service Service
}

func NewMetricsService(inner Service) *MetricsService {
	return &MetricsService{service: inner}
}

func (m *MetricsService) Increment(ctx context.Context, req models.IncrementPayload) (models.IncrementViewResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.Printf("Increment for VideoID %s took %v", req.VideoID, duration)
	}()

	return m.service.Increment(ctx, req)
}

func (m *MetricsService) Get(ctx context.Context, req models.ViewCountPayload) (models.ViewCountResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.Printf("Get for VideoID %s took %v", req.VideoID, duration)
	}()

	return m.service.Get(ctx, req)
}

func (m *MetricsService) GetTopVideos(ctx context.Context, req models.GetTopVideosPayload) (models.GetTopVideosResponse, error) {
	startTime := time.Now()

	resp, err := m.service.GetTopVideos(ctx, req)

	duration := time.Since(startTime)
	logMetrics("GetTopVideos", duration, err)
	return resp, err
}

func logMetrics(methodName string, duration time.Duration, err error) {
	// TODO prometheus?
	log.Printf("Metrics: %s took %v, Error: %v", methodName, duration, err)
}
