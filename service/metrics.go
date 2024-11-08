package service

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "method_duration_seconds",
			Help:    "Duration of each method in seconds.",
			Buckets: []float64{0.0001, 0.001, 0.005, 0.01, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"method"},
	)

	errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "method_errors_total",
			Help: "Total number of errors for each method.",
		},
		[]string{"method"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(errorCounter)
}

type MetricsService struct {
	service Service
}

func NewMetricsService(inner Service) *MetricsService {
	return &MetricsService{service: inner}
}

func (m *MetricsService) Increment(ctx context.Context, videoId string) (views int, increment int, err error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		logMetrics("Increment", duration, err)
	}()

	return m.service.Increment(ctx, videoId)
}

func (m *MetricsService) Get(ctx context.Context, videoId string) (views int, err error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		logMetrics("Get", duration, err)
	}()

	return m.service.Get(ctx, videoId)
}

func (m *MetricsService) GetTopVideos(ctx context.Context, page int, limit int) (topVideos []map[string]interface{}, err error) {
	startTime := time.Now()

	resp, err := m.service.GetTopVideos(ctx, page, limit)

	duration := time.Since(startTime)
	logMetrics("GetTopVideos", duration, err)
	return resp, err
}

func logMetrics(methodName string, duration time.Duration, err error) {
	// Record duration in Prometheus
	requestDuration.WithLabelValues(methodName).Observe(duration.Seconds())

	// Record an error if one occurred
	if err != nil {
		errorCounter.WithLabelValues(methodName).Inc()
	}
	log.Printf("%s took %v, Error: %v", methodName, duration, err)
}
