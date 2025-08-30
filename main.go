// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	//"os/user"
	"time"

	"github.com/edwinjordan/ZOGTest-Golang.git/config"
	"github.com/edwinjordan/ZOGTest-Golang.git/database"
	"github.com/edwinjordan/ZOGTest-Golang.git/domain"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/logging"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/metrics"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/repository/postgres"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/rest"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/rest/middleware"
	"github.com/edwinjordan/ZOGTest-Golang.git/service"
	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/edwinjordan/ZOGTest-Golang.git/docs"
)

func init() {
	config.LoadEnv()
}

// func adaptHandler(h http.Handler) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 		h.ServeHTTP(w, r)
// 	}
// }

func main() {
	config.SetupLogging()

	dbPool, err := database.SetupPgxPool()
	if err != nil {
		logging.LogError(context.Background(), err, "database_setup")
		os.Exit(1)
	}

	defer dbPool.Close()

	e := echo.New()
	e.HideBanner = true

	e.Logger.SetOutput(os.Stdout)
	e.Logger.SetLevel(0)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defer stop()

	appMetrics := metrics.NewMetrics()
	shutdown, err := config.ApplyInstrumentation(ctx, e, appMetrics)
	defer shutdown(ctx)
	e.Use(middleware.RequestIDMiddleware())
	e.Use(middleware.SlogLoggerMiddleware())
	e.Use(middleware.Cors())
	e.Use(middleware.SecurityHeadersMiddleware())
	e.Use(middleware.CompressionMiddleware())
	e.Use(middleware.RateLimitMiddleware(10.0, 20))
	e.Use(middleware.TimeoutMiddleware(30 * time.Second))
	//e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, domain.Response{
			Code:    200,
			Status:  "Succes",
			Message: "All is well!",
		})
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//e.Logger.Fatal(e.Start(":8080"))
	userRepo := postgres.NewUserRepository(dbPool, appMetrics)
	userService := service.NewUserService(userRepo)

	topicRepo := postgres.NewTopicRepository(dbPool)
	topicService := service.NewTopicService(topicRepo)

	newsRepo := postgres.NewNewsRepository(dbPool)
	newsService := service.NewNewsService(newsRepo)
	apiV1 := e.Group("/api/v1")
	usersGroup := apiV1.Group("")
	topicGroup := apiV1.Group("")
	newsGroup := apiV1.Group("")

	rest.NewUserHandler(usersGroup, userService)
	rest.NewTopicHandler(topicGroup, topicService)
	rest.NewNewsHandler(newsGroup, newsService)

	// Get host from environment variable, default to 127.0.0.1 if not set
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	// Get port from environment variable, default to 8000 if not set
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	// Server address and port to listen on
	serverAddr := fmt.Sprintf("%s:%s", host, port)

	go func() {
		logging.LogInfo(ctx, "Server starting", slog.String("address", serverAddr))
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			logging.LogError(ctx, err, "server_start")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logging.LogInfo(ctx, "Shutting down server gracefully...")
	if err := e.Shutdown(ctx); err != nil {
		logging.LogError(ctx, err, "server_shutdown")
	}
}
