package main

import (
	"youtubeviews/service"
	"youtubeviews/transport"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestIncrementHandler(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cacheCapacity := 20

	var svc service.Service
	{
		svc = service.CreateNewVideoService(redisClient, cacheCapacity)
		svc = service.NewLoggingService(svc)
		svc = service.NewMetricsService(svc)
	}

	transport := transport.NewHttpTransport(svc)

	req, err := http.NewRequest("POST", "/increment", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transport.Increment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetHandler(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cacheCapacity := 20

	var svc service.Service
	{
		svc = service.CreateNewVideoService(redisClient, cacheCapacity)
		svc = service.NewLoggingService(svc)
		svc = service.NewMetricsService(svc)
	}

	transport := transport.NewHttpTransport(svc)

	req, err := http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transport.Get)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetTopVideosHandler(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cacheCapacity := 20

	var svc service.Service
	{
		svc = service.CreateNewVideoService(redisClient, cacheCapacity)
		svc = service.NewLoggingService(svc)
		svc = service.NewMetricsService(svc)
	}

	transport := transport.NewHttpTransport(svc)

	req, err := http.NewRequest("GET", "/top-videos", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(transport.GetTopVideos)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
