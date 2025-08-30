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

func TestTopicCRUD_E2E(t *testing.T) {
	kit := NewTestKit(t)

	// Wire the routes and services
	//userRepo := postgres.NewUserRepository(kit.DB)
	//userSvc := service.NewUserService(userRepo)
	//rest.NewUserHandler(kit.Echo.Group("/api/v1"), userSvc)

	// Register Topic routes and services
	topicRepo := postgres.NewTopicRepository(kit.DB)
	topicSvc := service.NewTopicService(topicRepo)
	rest.NewTopicHandler(kit.Echo.Group("/api/v1"), topicSvc)

	// Now start the test server
	kit.Start(t)

	// Create
	createReq := domain.CreateTopicRequest{
		Name: "John Doe",
		//Slug: "john-doe",
		//Password: "Password1234",
	}
	type CreateType domain.ResponseSingleData[domain.Topic]
	cre, code := doRequest[CreateType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/topics",
		createReq,
	)
	require.Equal(t, http.StatusCreated, code)
	require.Equal(t, "success", cre.Status)
	topic := cre.Data
	require.NotEmpty(t, topic.ID)

	// Get
	type GetType domain.ResponseSingleData[domain.Topic]
	getE, code := doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/topics/%s", kit.BaseURL, topic.ID),
		nil,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, topic.ID, getE.Data.ID)

	// Update
	updPayload := domain.Topic{
		Name: "Jane Doe",
		Slug: "jane-doe",
	}
	type UpdType domain.ResponseSingleData[domain.Topic]
	updE, code := doRequest[UpdType](
		t, http.MethodPut,
		fmt.Sprintf("%s/api/v1/topics/%s", kit.BaseURL, topic.ID),
		updPayload,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, "Jane Doe", updE.Data.Name)

	// Delete
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/topics/%s", kit.BaseURL, topic.ID),
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
		fmt.Sprintf("%s/api/v1/topics/%s", kit.BaseURL, topic.ID),
		nil,
	)
	require.Equal(t, http.StatusNotFound, code)
	require.Equal(t, "error", errE.Status)
	require.Equal(t, "Topic not found", errE.Message)

	// Hard delete, since delete API uses soft delete
	_, err = kit.DB.Exec(context.Background(), "DELETE from topik where id = $1", topic.ID)
	require.NoError(t, err)
}
