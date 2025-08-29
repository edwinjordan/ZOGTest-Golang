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

func TestNewsService_CreateNews(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)

	newsService := service.NewNewsService(mockNewsRepo)

	ctx := context.Background()
	req := &domain.CreateNewsRequest{
		Title:   "Test News",
		Slug:    "test-news",
		Status:  "draft",
		Content: "This is a test news content.",
		Topic: []domain.NewsTopic{
			{
				TopicId: uuid.New().String(),
				NewsId:  uuid.New().String(),
			},
		},
	}
	expectedNews := &domain.News{
		ID:      uuid.New().String(),
		Title:   "Test News",
		Slug:    "test-news",
		Status:  "draft",
		Content: "This is a test news content.",
	}

	t.Run("Successfully creates a news", func(t *testing.T) {
		mockNewsRepo.On("CreateNews", mock.Anything, req).Return(expectedNews, nil).Once()

		news, err := newsService.CreateNews(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, news)
		assert.Equal(t, expectedNews.ID, news.ID)
		assert.Equal(t, expectedNews.Slug, news.Slug)
		//assert.Equal(t, expectedTop.Email, user.Email)

		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		repoErr := errors.New("database error")
		mockNewsRepo.On("CreateNews", mock.Anything, req).Return(nil, repoErr).Once()

		news, err := newsService.CreateNews(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Equal(t, repoErr, err)

		mockNewsRepo.AssertExpectations(t)
	})
}

func TestNewsService_GetNews(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	newsService := service.NewNewsService(mockNewsRepo)

	ctx := context.Background()
	newsID := uuid.New()
	expectedNews := &domain.News{
		ID:      newsID.String(),
		Title:   "Fetched Title",
		Slug:    "fetched-title",
		Status:  "published",
		Content: "This is fetched news content.",
		Topics: []domain.NewsTopicList{
			{
				ID:   uuid.New().String(),
				Name: "Tech",
			},
		},
	}

	t.Run("Successfully fetches a news", func(t *testing.T) {
		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(expectedNews, nil).Once()

		news, err := newsService.GetNews(ctx, newsID)

		assert.NoError(t, err)
		assert.NotNil(t, news)
		assert.Equal(t, expectedNews.ID, news.ID)
		assert.Equal(t, expectedNews.Title, news.Title)
		assert.Equal(t, expectedNews.Slug, news.Slug)
		assert.Equal(t, expectedNews.Status, news.Status)
		assert.Equal(t, expectedNews.Content, news.Content)

		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		repoErr := errors.New("network error")
		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(nil, repoErr).Once()

		news, err := newsService.GetNews(ctx, newsID)

		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Equal(t, repoErr, err)

		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns nil when topic not found in repository", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(nil, nil).Once()

		news, err := newsService.GetNews(ctx, newsID)

		assert.NoError(t, err)
		assert.Nil(t, news)

		mockNewsRepo.AssertExpectations(t)
	})
}

func TestNewsService_UpdateNews(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	newsService := service.NewNewsService(mockNewsRepo)

	ctx := context.Background()
	newsID := uuid.New()
	existingNews := &domain.News{
		ID:      newsID.String(),
		Title:   "Old Name",
		Slug:    "old-name",
		Status:  "draft",
		Content: "Old content",
		// Topic: []domain.NewsTopic{
		// 	{
		// 		TopicId: uuid.New().String(),
		// 		NewsId:  newsID.String(),
		// 	},
		// },
	}
	updateReq := &domain.News{
		Title:   "New Name",
		Slug:    "new-name",
		Status:  "published",
		Content: "New content",
	}

	t.Run("Successfully updates a news", func(t *testing.T) {
		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(existingNews, nil).Once()

		expectedUpdatedNews := &domain.News{
			ID:      newsID.String(),
			Title:   updateReq.Title,
			Slug:    updateReq.Slug,
			Status:  updateReq.Status,
			Content: updateReq.Content,
		}
		mockNewsRepo.On("UpdateNews", mock.Anything, newsID, expectedUpdatedNews).Return(expectedUpdatedNews, nil).Once()

		user, err := newsService.UpdateNews(ctx, newsID, updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUpdatedNews.Title, user.Title)
		assert.Equal(t, expectedUpdatedNews.Slug, user.Slug)
		assert.Equal(t, expectedUpdatedNews.Status, user.Status)
		assert.Equal(t, expectedUpdatedNews.Content, user.Content)

		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrNewsNotFound if user does not exist", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(nil, nil).Once()

		topic, err := newsService.UpdateNews(ctx, newsID, updateReq)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, topic)

		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetNews fails", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		repoErr := errors.New("get news repo error")
		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(nil, repoErr).Once()

		user, err := newsService.UpdateNews(ctx, newsID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, repoErr, err)

		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if UpdateNews fails", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(existingNews, nil).Once()

		repoErr := errors.New("update news repo error")
		expectedUpdatedNews := &domain.News{
			ID:      newsID.String(),
			Title:   updateReq.Title,
			Slug:    updateReq.Slug,
			Status:  updateReq.Status,
			Content: updateReq.Content,
		}
		mockNewsRepo.On("UpdateNews", mock.Anything, newsID, expectedUpdatedNews).Return(nil, repoErr).Once()

		news, err := newsService.UpdateNews(ctx, newsID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Equal(t, repoErr, err)

		mockNewsRepo.AssertExpectations(t)
	})
}

func TestNewsService_DeleteNews(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	newsService := service.NewNewsService(mockNewsRepo)

	ctx := context.Background()
	newsID := uuid.New()
	existingNews := &domain.News{
		ID:      newsID.String(),
		Title:   "User to delete",
		Slug:    "user-to-delete",
		Status:  "draft",
		Content: "Content to delete",
	}

	t.Run("Successfully deletes a news", func(t *testing.T) {
		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(existingNews, nil).Once()
		mockNewsRepo.On("DeleteNews", mock.Anything, newsID).Return(nil).Once()

		err := newsService.DeleteNews(ctx, newsID)

		assert.NoError(t, err)
		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrNewsNotFound if user does not exist", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(nil, nil).Once()

		err := newsService.DeleteNews(ctx, newsID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetNews fails", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		repoErr := errors.New("get news repo error during delete")
		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(nil, repoErr).Once()

		err := newsService.DeleteNews(ctx, newsID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns error if DeleteNews fails", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		mockNewsRepo.On("GetNews", mock.Anything, newsID).Return(existingNews, nil).Once()
		repoErr := errors.New("delete news repo error")
		mockNewsRepo.On("DeleteNews", mock.Anything, newsID).Return(repoErr).Once()

		err := newsService.DeleteNews(ctx, newsID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockNewsRepo.AssertExpectations(t)
	})
}

func TestNewsService_GetNewsList(t *testing.T) {
	mockNewsRepo := new(mocks.NewsRepository)
	newsService := service.NewNewsService(mockNewsRepo)

	ctx := context.Background()
	filter := &domain.NewsFilter{
		Search: "test",
	}
	expectedNews := []domain.News{
		{ID: uuid.New().String(), Title: "Test News One", Slug: "test-news-one"},
		{ID: uuid.New().String(), Title: "Another Test News", Slug: "another-test-news"},
	}

	t.Run("Successfully fetches news list", func(t *testing.T) {
		mockNewsRepo.On("GetNewsList", mock.Anything, filter).Return(expectedNews, nil).Once()

		news, err := newsService.GetNewsList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, news)
		assert.Len(t, news, 2)
		assert.Equal(t, expectedNews[0].Title, news[0].Title)

		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns empty list when no news found", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		mockNewsRepo.On("GetNewsList", mock.Anything, filter).Return([]domain.News{}, nil).Once()

		news, err := newsService.GetNewsList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, news)
		assert.Len(t, news, 0)

		mockNewsRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockNewsRepo = new(mocks.NewsRepository)
		newsService = service.NewNewsService(mockNewsRepo)

		repoErr := errors.New("get news list database error")
		mockNewsRepo.On("GetNewsList", mock.Anything, filter).Return(nil, repoErr).Once()

		news, err := newsService.GetNewsList(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Equal(t, repoErr, err)

		mockNewsRepo.AssertExpectations(t)
	})
}
