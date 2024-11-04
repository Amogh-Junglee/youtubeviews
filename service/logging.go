package service

import (
	models "YoutubeViews/models"
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
	log.Printf("Started Increment for VideoID: %s", req.VideoID)
	log.Println("Request: ", req)

	resp, err := l.service.Increment(ctx, req)
	if err != nil {
		log.Printf("Error Incrementing VideoID %s: %v", req.VideoID, err)
		return models.IncrementViewResponse{}, err
	}

	log.Println("Response: ", resp)
	log.Printf("Completed Increment for VideoID: %s, Views: %d, Increment: %d", req.VideoID, resp.Views, resp.Increment)
	return resp, nil
}

func (l *LoggingService) Get(ctx context.Context, req models.ViewCountPayload) (models.ViewCountResponse, error) {
	log.Printf("Started Get for VideoID: %s", req.VideoID)
	log.Println("Request: ", req)

	resp, err := l.service.Get(ctx, req)
	if err != nil {
		log.Printf("Error getting view count for VideoID %s: %v", req.VideoID, err)
		return resp, err
	}

	log.Println("Response: ", resp)
	log.Printf("Completed Get for VideoID: %s, Views: %d", req.VideoID, resp.GetViews())
	return resp, nil
}

func (l *LoggingService) GetTopVideos(ctx context.Context, req models.GetTopVideosPayload) (models.GetTopVideosResponse, error) {
	startTime := time.Now()
	log.Println("Request: ", req)

	// Call the next service in the chain
	resp, err := l.service.GetTopVideos(ctx, req)

	// Log request details
	log.Println("Response: ", resp)
	log.Printf("Request: GetTopVideos, Page: %d, Limit: %d, TimeTaken: %v, Error: %v", req.Page, req.Limit, time.Since(startTime), err)
	return resp, err
}
