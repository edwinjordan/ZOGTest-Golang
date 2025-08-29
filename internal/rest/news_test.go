package rest_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/repository/postgres"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/rest"
	"github.com/edwinjordan/ZOGTest-Golang.git/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewsCRUD_E2E(t *testing.T) {
	kit := NewTestKit(t)

	// Wire the routes and services
	//userRepo := postgres.NewUserRepository(kit.DB)
	//userSvc := service.NewUserService(userRepo)
	//rest.NewUserHandler(kit.Echo.Group("/api/v1"), userSvc)

	// Register User routes and services
	newsRepo := postgres.NewNewsRepository(kit.DB)
	newsSvc := service.NewNewsService(newsRepo)
	rest.NewNewsHandler(kit.Echo.Group("/api/v1"), newsSvc)

	// Now start the test server
	kit.Start(t)

	// Create
	createReq := domain.CreateNewsRequest{
		Title:   "Test News",
		Status:  "draft",
		Content: "This is a test news content.",
		Topic: []domain.NewsTopic{
			{
				TopicId: uuid.New().String(),
				NewsId:  uuid.New().String(),
			},
		},
		//Slug: "john-doe",
		//Password: "Password1234",
	}
	type CreateType domain.ResponseSingleData[domain.News]
	cre, code := doRequest[CreateType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/news",
		createReq,
	)
	require.Equal(t, http.StatusCreated, code)
	require.Equal(t, "success", cre.Status)
	news := cre.Data
	require.NotEmpty(t, news.ID)

	// Get
	type GetType domain.ResponseSingleData[domain.News]
	getE, code := doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID),
		nil,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, news.ID, getE.Data.ID)

	// // Update
	// updPayload := domain.UpdateNewsRequest{
	// 	Title:   "Updated Test News",
	// 	Status:  "published",
	// 	Content: "This is the updated content.",
	// 	Topic: []domain.NewsTopic{
	// 		{
	// 			TopicId: "5826280c-41f3-4d2b-a09d-adc85f07d8ac",
	// 			//NewsId:  uuid.New().String(),
	// 		},
	// 	},
	// }
	// type UpdType domain.ResponseSingleData[domain.News]
	// updE, code := doRequest[UpdType](
	// 	t, http.MethodPut,
	// 	fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID),
	// 	updPayload,
	// )
	// require.Equal(t, http.StatusOK, code)
	// require.Equal(t, "Updated Test News", updE.Data.Title)

	// Delete
	// req, err := http.NewRequest(
	// 	http.MethodDelete,
	// 	fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID),
	// 	nil,
	// )
	// require.NoError(t, err)
	// resp, err := http.DefaultClient.Do(req)
	// require.NoError(t, err)
	// require.Equal(t, http.StatusNoContent, resp.StatusCode)
	// resp.Body.Close()

	// Get after delete
	// type ErrType domain.ResponseSingleData[domain.Empty]
	// errE, code := doRequest[ErrType](
	// 	t, http.MethodGet,
	// 	fmt.Sprintf("%s/api/v1/news/%s", kit.BaseURL, news.ID),
	// 	nil,
	// )
	// require.Equal(t, http.StatusNotFound, code)
	// require.Equal(t, "error", errE.Status)
	// require.Equal(t, "News not found", errE.Message)

	// // Hard delete, since delete API uses soft delete
	// _, err = kit.DB.Exec(context.Background(), "DELETE from news where id = $1", news.ID)
	// require.NoError(t, err)
}
