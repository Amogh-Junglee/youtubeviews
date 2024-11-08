package main

import (
	service "youtubeviews/service"
	transport "youtubeviews/transport"

	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

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

	http.HandleFunc("/increment", transport.Increment)
	http.HandleFunc("/get", transport.Get)
	http.HandleFunc("/top-videos", transport.GetTopVideos)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
