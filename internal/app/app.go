package app

import (
	"context"
	"fmt"
	"github.com/core-go/health"
	es "github.com/core-go/health/elasticsearch/v8"
	"github.com/core-go/log"
	"github.com/elastic/go-elasticsearch/v8"

	"go-service/internal/handler"
	"go-service/internal/service"
)

type ApplicationContext struct {
	Health *health.Handler
	User   *handler.UserHandler
}

func NewApp(ctx context.Context, config Config) (*ApplicationContext, error) {
	log.Initialize(config.Log)

	cfg := elasticsearch.Config{Addresses: []string{config.ElasticSearch.Url}}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Error(ctx, "Cannot connect to elasticSearch. Error: "+err.Error())
		return nil, err
	}

	res, err := client.Info()
	if err != nil {
		log.Error(ctx, "Elastic server Error: " + err.Error())
		return nil, err
	}
	fmt.Println("Elastic server response: ", res)

	userService := service.NewUserService(client)
	userHandler := handler.NewUserHandler(userService)

	elasticSearchChecker := es.NewHealthChecker(client)
	healthHandler := health.NewHandler(elasticSearchChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
