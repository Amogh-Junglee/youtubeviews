package service

import (
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

func (l *LoggingService) Increment(ctx context.Context, videoId string) (views int, increment int, err error) {
	log.Printf("LoggingService.Increment: Started Increment for VideoID: %s", videoId)

	views, inc, err := l.service.Increment(ctx, videoId)
	if err != nil {
		log.Printf("LoggingService.Increment: Error Incrementing VideoID %s: %v", videoId, err)
		return -1, -1, err
	}

	log.Println("LoggingService.Increment: Response: ", inc, " ", views)
	log.Printf("LoggingService.Increment: Completed Increment for VideoID: %s, Views: %d, Increment: %d", videoId, views, inc)
	return views, inc, nil
}

func (l *LoggingService) Get(ctx context.Context, videoId string) (views int, err error) {
	log.Printf("LoggingService.Get: Started Get for VideoID: %s", videoId)

	views, err = l.service.Get(ctx, videoId)
	if err != nil {
		log.Printf("LoggingService.Get: Error getting view count for VideoID %s: %v", videoId, err)
		return views, err
	}

	log.Println("LoggingService.Get: Response: ", views)
	log.Printf("LoggingService.Get: Completed Get for VideoID: %s, Views: %d", videoId, views)
	return views, nil
}

func (l *LoggingService) GetTopVideos(ctx context.Context, page int, limit int) (topVideos []map[string]interface{}, err error) {
	startTime := time.Now()
	log.Println("LoggingService.GetTopVideos: Request: ", page, " ", limit)

	// Call the next service in the chain
	resp, err := l.service.GetTopVideos(ctx, page, limit)

	// Log request details
	log.Println("LoggingService.GetTopVideos: Response: ", resp)
	log.Printf("LoggingService.GetTopVideos: Request: GetTopVideos, Page: %d, Limit: %d, TimeTaken: %v, Error: %v", page, limit, time.Since(startTime), err)
	return resp, err
}
