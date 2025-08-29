package rest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/repository/postgres"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/rest"
	"github.com/edwinjordan/ZOGTest-Golang.git/service"
	"github.com/stretchr/testify/require"
)

func TestNewsCRUD_E2E(t *testing.T) {
	kit := NewTestKit(t)

	// Wire the routes and services
	//userRepo := postgres.NewUserRepository(kit.DB)
	//userSvc := service.NewUserService(userRepo)
	//rest.NewUserHandler(kit.Echo.Group("/api/v1"), userSvc)
	topicRepo := postgres.NewTopicRepository(kit.DB)
	topicSvc := service.NewTopicService(topicRepo)
	rest.NewTopicHandler(kit.Echo.Group("/api/v1"), topicSvc)
	// Register Topic routes and services
	newsRepo := postgres.NewNewsRepository(kit.DB)
	newsSvc := service.NewNewsService(newsRepo)
	rest.NewNewsHandler(kit.Echo.Group("/api/v1"), newsSvc)

	// Now start the test server
	kit.Start(t)

	// Create Topik
	createReq2 := domain.CreateTopicRequest{
		Name: "Test Topic",
		//Slug: "john-doe",
		//Password: "Password1234",
	}
	type CreateType2 domain.ResponseSingleData[domain.Topic]
	cre2, code2 := doRequest[CreateType2](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/topics",
		createReq2,
	)
	require.Equal(t, http.StatusCreated, code2)
	require.Equal(t, "success", cre2.Status)
	topic := cre2.Data
	require.NotEmpty(t, topic.ID)

	// Create News
	createReq := domain.CreateNewsRequest{
		Title:   "Breaking News",
		Slug:    "breaking-news",
		Status:  "published",
		Content: "This is the content of the breaking news.",
		Topic: []domain.NewsTopic{
			{
				TopicId: topic.ID,
			},
		},
	}
	type CreateType domain.ResponseSingleData[domain.Topic]
	cre, code := doRequest[CreateType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/news",
		createReq,
	)
	require.Equal(t, http.StatusCreated, code)
	require.Equal(t, "success", cre.Status)
	news := cre.Data
	require.NotEmpty(t, news.ID)

	// // Get
	type GetType domain.ResponseSingleData[domain.News]
	getE, code := doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID),
		nil,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, news.ID, getE.Data.ID)

	// // Update

	updPayload := domain.News{
		Title:   "Updated News Title",
		Status:  "draft",
		Content: "This is the updated content of the news.",
		Topics: []domain.NewsTopic{
			{
				TopicId: topic.ID,
			},
		},
	}
	type UpdType domain.ResponseSingleData[domain.News]
	updE, code := doRequest[UpdType](
		t, http.MethodPut,
		fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID),
		updPayload,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, "Updated News Title", updE.Data.Title)

	// // Delete
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID),
		nil,
	)

	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Get after delete
	type ErrType domain.ResponseSingleData[domain.Empty]
	errE, code := doRequest[ErrType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID),
		nil,
	)
	require.Equal(t, http.StatusNotFound, code)
	require.Equal(t, "error", errE.Status)
	require.Equal(t, "News not found", errE.Message)

	// Hard delete, since delete API uses soft delete
	_, err = kit.DB.Exec(context.Background(), "DELETE from news where id = $1", news.ID)
	require.NoError(t, err)
}
