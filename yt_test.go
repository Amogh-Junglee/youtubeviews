package main

import (
	models "YoutubeViews/models"
	service "YoutubeViews/service"
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/redis/go-redis/v9"
)

const capacity = 20

func setupRedis() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return redisClient, nil
}

func cleanupRedis(redisClient *redis.Client) {
	redisClient.FlushAll(context.Background())
	redisClient.Close()
}

func mapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !reflect.DeepEqual(v, bv) {
			return false
		}
	}
	return true
}

func slicesOfMapsEqual(a, b []map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !mapsEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func debugDeepEqual(result, want interface{}) {
	if reflect.TypeOf(result) != reflect.TypeOf(want) {
		fmt.Printf("Type mismatch: result is %T, want is %T\n", result, want)
		return
	}

	resultVal := reflect.ValueOf(result)
	wantVal := reflect.ValueOf(want)

	if resultVal.Len() != wantVal.Len() {
		fmt.Printf("Length mismatch: result has length %d, want has length %d\n", resultVal.Len(), wantVal.Len())
		return
	}

	for i := 0; i < resultVal.Len(); i++ {
		if !reflect.DeepEqual(resultVal.Index(i).Interface(), wantVal.Index(i).Interface()) {
			fmt.Printf("Element mismatch at index %d: result[%d] = %+v, want[%d] = %+v\n", i, i, resultVal.Index(i).Interface(), i, wantVal.Index(i).Interface())
		}
	}
}

func TestIncrementAndGetWithCache(t *testing.T) {
	redisClient, err := setupRedis()
	if err != nil {
		t.Fatalf("Could not connect to Redis: %v", err)
	}
	defer cleanupRedis(redisClient)

	var svc service.Service
	{
		svc = service.CreateNewVideoService(redisClient, capacity)
		svc = service.NewLoggingService(svc)
		svc = service.NewMetricsService(svc)
	}

	tests := []struct {
		videoID   string
		wantViews int
	}{
		{"video1", 1},
		{"video1", 2},
		{"video2", 1},
	}

	for _, tt := range tests {
		t.Run(tt.videoID, func(t *testing.T) {
			result, err := svc.Increment(context.Background(), models.IncrementPayload{VideoID: tt.videoID})
			if err != nil {
				t.Errorf("Increment() error = %v", err)
			}

			_, getError := svc.Get(context.Background(), models.ViewCountPayload{VideoID: tt.videoID})
			if getError != nil {
				t.Errorf("Get() error = %v", getError)
			}

			if result.Views != tt.wantViews {
				t.Errorf("Get() = %v, want %v", result.Views, tt.wantViews)
			}
		})
	}
	fmt.Println("TestIncrementAndGetWithCache complete")
}

func TestGetTopVideosWithCache(t *testing.T) {
	redisClient, err := setupRedis()
	if err != nil {
		t.Fatalf("Could not connect to Redis: %v", err)
	}
	defer cleanupRedis(redisClient)

	var svc service.Service
	{
		svc = service.CreateNewVideoService(redisClient, capacity)
		svc = service.NewLoggingService(svc)
		svc = service.NewMetricsService(svc)
	}

	// Prerequisite data
	videoIDs := []string{"video1", "video2", "video3"}
	for _, videoID := range videoIDs {
		for i := 0; i < 3; i++ {
			svc.Increment(context.Background(), models.IncrementPayload{VideoID: videoID})
		}
	}

	fmt.Println("Prerequisite complete")

	tests := []struct {
		name       string
		wantVideos []map[string]interface{}
	}{
		{
			name: "GetTopVideos",
			wantVideos: []map[string]interface{}{
				{"videoId": "video3", "views": float64(3)},
				{"videoId": "video2", "views": float64(3)},
				{"videoId": "video1", "views": float64(3)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetTopVideos(context.Background(), models.GetTopVideosPayload{Limit: 3, Page: 1})
			if err != nil {
				t.Fatalf("GetTopVideos() error = %v", err)
			}

			if !slicesOfMapsEqual(result.TopVideos, tt.wantVideos) {
				t.Errorf("GetTopVideos() = %+v, want %+v", result.TopVideos, tt.wantVideos)
				// debugDeepEqual(result.TopVideos, tt.wantVideos)
			} else {
				fmt.Println("GetTopVideos() matches expected output.")
			}
		})
	}
}
