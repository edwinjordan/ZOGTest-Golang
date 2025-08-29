package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/service"
	"github.com/edwinjordan/ZOGTest-Golang.git/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTopicService_CreateTopic(t *testing.T) {
	mockTopicRepo := new(mocks.TopicRepository)

	topicService := service.NewTopicService(mockTopicRepo)

	ctx := context.Background()
	req := &domain.CreateTopicRequest{
		Name: "Test Topic",
		Slug: "test-topic",
	}
	expectedTopic := &domain.Topic{
		ID:   uuid.New().String(),
		Name: "Test Topic",
		Slug: "test-topic",
	}

	t.Run("Successfully creates a topic", func(t *testing.T) {
		mockTopicRepo.On("CreateTopic", mock.Anything, req).Return(expectedTopic, nil).Once()

		topic, err := topicService.CreateTopic(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, topic)
		assert.Equal(t, expectedTopic.ID, topic.ID)
		assert.Equal(t, expectedTopic.Slug, topic.Slug)
		//assert.Equal(t, expectedTop.Email, user.Email)

		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		repoErr := errors.New("database error")
		mockTopicRepo.On("CreateTopic", mock.Anything, req).Return(nil, repoErr).Once()

		topic, err := topicService.CreateTopic(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, topic)
		assert.Equal(t, repoErr, err)

		mockTopicRepo.AssertExpectations(t)
	})
}

func TestTopicService_GetTopic(t *testing.T) {
	mockTopicRepo := new(mocks.TopicRepository)
	topicService := service.NewTopicService(mockTopicRepo)

	ctx := context.Background()
	topicID := uuid.New()
	expectedTopic := &domain.Topic{
		ID:   topicID.String(),
		Name: "Fetched User",
		Slug: "fetched-user",
	}

	t.Run("Successfully fetches a topic", func(t *testing.T) {
		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(expectedTopic, nil).Once()

		topic, err := topicService.GetTopic(ctx, topicID)

		assert.NoError(t, err)
		assert.NotNil(t, topic)
		assert.Equal(t, expectedTopic.ID, topic.ID)
		assert.Equal(t, expectedTopic.Name, topic.Name)

		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		repoErr := errors.New("network error")
		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(nil, repoErr).Once()

		user, err := topicService.GetTopic(ctx, topicID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, repoErr, err)

		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns nil when topic not found in repository", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(nil, nil).Once()

		user, err := topicService.GetTopic(ctx, topicID)

		assert.NoError(t, err)
		assert.Nil(t, user)

		mockTopicRepo.AssertExpectations(t)
	})
}

func TestTopicService_UpdateTopic(t *testing.T) {
	mockTopicRepo := new(mocks.TopicRepository)
	topicService := service.NewTopicService(mockTopicRepo)

	ctx := context.Background()
	topicID := uuid.New()
	existingTopic := &domain.Topic{
		ID:   topicID.String(),
		Name: "Old Name",
		Slug: "old-name",
	}
	updateReq := &domain.Topic{
		Name: "New Name",
		Slug: "new-name",
	}

	t.Run("Successfully updates a topic", func(t *testing.T) {
		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(existingTopic, nil).Once()

		expectedUpdatedTopic := &domain.Topic{
			ID:   topicID.String(),
			Name: updateReq.Name,
			Slug: updateReq.Slug,
		}
		mockTopicRepo.On("UpdateTopic", mock.Anything, topicID, expectedUpdatedTopic).Return(expectedUpdatedTopic, nil).Once()

		user, err := topicService.UpdateTopic(ctx, topicID, updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUpdatedTopic.Name, user.Name)
		assert.Equal(t, expectedUpdatedTopic.Slug, user.Slug)

		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrTopicNotFound if user does not exist", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(nil, nil).Once()

		topic, err := topicService.UpdateTopic(ctx, topicID, updateReq)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, topic)

		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetTopic fails", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		repoErr := errors.New("get topic repo error")
		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(nil, repoErr).Once()

		user, err := topicService.UpdateTopic(ctx, topicID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, repoErr, err)

		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns error if UpdateTopic fails", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(existingTopic, nil).Once()

		repoErr := errors.New("update topic repo error")
		expectedUpdatedTopic := &domain.Topic{
			ID:   topicID.String(),
			Name: updateReq.Name,
			Slug: updateReq.Slug,
		}
		mockTopicRepo.On("UpdateTopic", mock.Anything, topicID, expectedUpdatedTopic).Return(nil, repoErr).Once()

		user, err := topicService.UpdateTopic(ctx, topicID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, repoErr, err)

		mockTopicRepo.AssertExpectations(t)
	})
}

func TestTopicService_DeleteTopic(t *testing.T) {
	mockTopicRepo := new(mocks.TopicRepository)
	topicService := service.NewTopicService(mockTopicRepo)

	ctx := context.Background()
	topicID := uuid.New()
	existingTopic := &domain.Topic{
		ID:   topicID.String(),
		Name: "User to delete",
		Slug: "user-to-delete",
	}

	t.Run("Successfully deletes a topic", func(t *testing.T) {
		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(existingTopic, nil).Once()
		mockTopicRepo.On("DeleteTopic", mock.Anything, topicID).Return(nil).Once()

		err := topicService.DeleteTopic(ctx, topicID)

		assert.NoError(t, err)
		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrTopicNotFound if user does not exist", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(nil, nil).Once()

		err := topicService.DeleteTopic(ctx, topicID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetTopic fails", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		repoErr := errors.New("get topic repo error during delete")
		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(nil, repoErr).Once()

		err := topicService.DeleteTopic(ctx, topicID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns error if DeleteTopic fails", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		mockTopicRepo.On("GetTopic", mock.Anything, topicID).Return(existingTopic, nil).Once()
		repoErr := errors.New("delete topic repo error")
		mockTopicRepo.On("DeleteTopic", mock.Anything, topicID).Return(repoErr).Once()

		err := topicService.DeleteTopic(ctx, topicID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockTopicRepo.AssertExpectations(t)
	})
}

func TestTopicService_GetTopicList(t *testing.T) {
	mockTopicRepo := new(mocks.TopicRepository)
	topicService := service.NewTopicService(mockTopicRepo)

	ctx := context.Background()
	filter := &domain.TopicFilter{
		Search: "test",
	}
	expectedTopic := []domain.Topic{
		{ID: uuid.New().String(), Name: "Test Topic One"},
		{ID: uuid.New().String(), Name: "Another Test Topic"},
	}

	t.Run("Successfully fetches topic list", func(t *testing.T) {
		mockTopicRepo.On("GetTopicList", mock.Anything, filter).Return(expectedTopic, nil).Once()

		topic, err := topicService.GetTopicList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, topic)
		assert.Len(t, topic, 2)
		assert.Equal(t, expectedTopic[0].Name, topic[0].Name)

		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns empty list when no topic found", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		mockTopicRepo.On("GetTopicList", mock.Anything, filter).Return([]domain.Topic{}, nil).Once()

		topics, err := topicService.GetTopicList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, topics)
		assert.Len(t, topics, 0)

		mockTopicRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockTopicRepo = new(mocks.TopicRepository)
		topicService = service.NewTopicService(mockTopicRepo)

		repoErr := errors.New("get topic list database error")
		mockTopicRepo.On("GetTopicList", mock.Anything, filter).Return(nil, repoErr).Once()

		topics, err := topicService.GetTopicList(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, topics)
		assert.Equal(t, repoErr, err)

		mockTopicRepo.AssertExpectations(t)
	})
}
