package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/adapter/database"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/handler/customer/delete_account"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/handler/health"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-customer-management/internal/shared/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	customer_repository "github.com/jfelipearaujo-org/ms-customer-management/internal/repository/customer"
	delete_request_repository "github.com/jfelipearaujo-org/ms-customer-management/internal/repository/delete_request"
	customer_delete_account_svc "github.com/jfelipearaujo-org/ms-customer-management/internal/service/customer/delete_account"
)

type Server struct {
	Config          *environment.Config
	DatabaseService database.DatabaseService

	Dependency Dependency
}

func NewServer(config *environment.Config) *Server {
	ctx := context.Background()

	cloudConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	if config.CloudConfig.IsBaseEndpointSet() {
		cloudConfig.BaseEndpoint = aws.String(config.CloudConfig.BaseEndpoint)
	}

	databaseService := database.NewDatabase(config)

	timeProvider := time_provider.NewTimeProvider(time.Now)

	customer_repository := customer_repository.NewRepository(databaseService.GetInstance())
	delete_request_repository := delete_request_repository.NewRepository(databaseService.GetInstance())

	return &Server{
		Config:          config,
		DatabaseService: databaseService,
		Dependency: Dependency{
			TimeProvider: timeProvider,

			CustomerRepository: customer_repository,
			CustomerService:    customer_delete_account_svc.NewService(customer_repository, delete_request_repository),
		},
	}
}

func (s *Server) GetHttpServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.ApiConfig.Port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(logger.Middleware())
	e.Use(middleware.Recover())

	s.registerHealthCheck(e)

	group := e.Group(fmt.Sprintf("/api/%s", s.Config.ApiConfig.ApiVersion))

	s.registerCustomerHandlers(group)

	return e
}

func (server *Server) registerHealthCheck(e *echo.Echo) {
	healthHandler := health.NewHandler(server.DatabaseService)

	e.GET("/health", healthHandler.Handle)
}

func (s *Server) registerCustomerHandlers(e *echo.Group) {
	customerHandler := delete_account.NewHandler(s.Dependency.CustomerService)

	e.POST("/customers/:id/delete-account", customerHandler.Handle)
}
