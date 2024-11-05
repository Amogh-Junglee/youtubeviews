package service

import (
	"context"
	"testing"

	"youtubeviews/db/mocks"

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
	req := "123"

	mockRepo.EXPECT().Increment(ctx, req).Return(1,1, nil)

	views, increment, err := videoService.Increment(ctx, req)
	assert.NoError(t, err)
	assert.Greater(t, views, 0)
	assert.True(t, increment == 1)
}

func TestVideoService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDbRepo(ctrl)
	videoService := &VideoService{repo: mockRepo}

	ctx := context.Background()
	req := "123"

	mockRepo.EXPECT().Get(ctx, req).Return(100, nil)

	views, err := videoService.Get(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 100, views)
}

func TestVideoService_GetTopVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDbRepo(ctrl)
	videoService := &VideoService{repo: mockRepo}

	ctx := context.Background()
	page,limit := 1,10

	mockRepo.EXPECT().GetTopVideos(ctx, page, limit).Return([]map[string]interface{}{{"ID": "123", "Views": 1000}}, nil)

	topVideos, err := videoService.GetTopVideos(ctx, page, limit)
	assert.NoError(t, err)
	assert.Len(t, topVideos, 1)
	topVideo := topVideos[0]
	id, idOk := topVideo["ID"].(string)
	views, viewsOk := topVideo["Views"].(int)

	assert.True(t, idOk, "ID should be a string")
	assert.True(t, viewsOk, "Views should be a float64")
	assert.Equal(t, "123", id)
	assert.Equal(t, int(1000), views)
}
