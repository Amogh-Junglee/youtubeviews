package service

import (
	"context"
	"testing"

	"youtubeviews/db/mocks"
	"youtubeviews/models"

	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewVideoService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := redis.NewClient(&redis.Options{})
	videoService := CreateNewVideoService(mockRedisClient, 100)

	assert.NotNil(t, videoService)
}

func TestVideoService_Increment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDbRepo(ctrl)
	videoService := &VideoService{repo: mockRepo}

	ctx := context.Background()
	req := models.IncrementPayload{VideoID: "123"}

	mockRepo.EXPECT().Increment(ctx, req).Return(models.IncrementViewResponse{Views: 1, Increment: 1}, nil)

	resp, err := videoService.Increment(ctx, req)
	assert.NoError(t, err)
	assert.True(t, resp.Increment == 1)
}

func TestVideoService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDbRepo(ctrl)
	videoService := &VideoService{repo: mockRepo}

	ctx := context.Background()
	req := models.ViewCountPayload{VideoID: "123"}

	mockRepo.EXPECT().Get(ctx, req).Return(models.ViewCountResponse{Views: 100}, nil)

	resp, err := videoService.Get(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 100, resp.Views)
}

func TestVideoService_GetTopVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDbRepo(ctrl)
	videoService := &VideoService{repo: mockRepo}

	ctx := context.Background()
	req := models.GetTopVideosPayload{Limit: 10}

	mockRepo.EXPECT().GetTopVideos(ctx, req).Return(models.GetTopVideosResponse{TopVideos: []map[string]interface{}{{"ID": "123", "Views": 1000}}}, nil)

	resp, err := videoService.GetTopVideos(ctx, req)
	assert.NoError(t, err)
	assert.Len(t, resp.TopVideos, 1)
	topVideo := resp.TopVideos[0]
	id, idOk := topVideo["ID"].(string)
	views, viewsOk := topVideo["Views"].(int)

	assert.True(t, idOk, "ID should be a string")
	assert.True(t, viewsOk, "Views should be a float64")
	assert.Equal(t, "123", id)
	assert.Equal(t, int(1000), views)
}
