package service

import (
	"YoutubeViews/models"
	"context"
	"log"
	"time"
)

type LoggingService struct {
	service Service
}

func NewLoggingService(inner Service) *LoggingService {
	return &LoggingService{service: inner}
}

func (l *LoggingService) Increment(ctx context.Context, req models.IncrementPayload) (models.IncrementViewResponse, error) {
	log.Printf("LoggingService.Increment: Started Increment for VideoID: %s", req.VideoID)
	log.Println("LoggingService.Increment: Request: ", req)

	resp, err := l.service.Increment(ctx, req)
	if err != nil {
		log.Printf("LoggingService.Increment: Error Incrementing VideoID %s: %v", req.VideoID, err)
		return models.IncrementViewResponse{}, err
	}

	log.Println("LoggingService.Increment: Response: ", resp)
	log.Printf("LoggingService.Increment: Completed Increment for VideoID: %s, Views: %d, Increment: %d", req.VideoID, resp.Views, resp.Increment)
	return resp, nil
}

func (l *LoggingService) Get(ctx context.Context, req models.ViewCountPayload) (models.ViewCountResponse, error) {
	log.Printf("LoggingService.Get: Started Get for VideoID: %s", req.VideoID)
	log.Println("LoggingService.Get: Request: ", req)

	resp, err := l.service.Get(ctx, req)
	if err != nil {
		log.Printf("LoggingService.Get: Error getting view count for VideoID %s: %v", req.VideoID, err)
		return resp, err
	}

	log.Println("LoggingService.Get: Response: ", resp)
	log.Printf("LoggingService.Get: Completed Get for VideoID: %s, Views: %d", req.VideoID, resp.GetViews())
	return resp, nil
}

func (l *LoggingService) GetTopVideos(ctx context.Context, req models.GetTopVideosPayload) (models.GetTopVideosResponse, error) {
	startTime := time.Now()
	log.Println("LoggingService.GetTopVideos: Request: ", req)

	// Call the next service in the chain
	resp, err := l.service.GetTopVideos(ctx, req)

	// Log request details
	log.Println("LoggingService.GetTopVideos: Response: ", resp)
	log.Printf("LoggingService.GetTopVideos: Request: GetTopVideos, Page: %d, Limit: %d, TimeTaken: %v, Error: %v", req.Page, req.Limit, time.Since(startTime), err)
	return resp, err
}
