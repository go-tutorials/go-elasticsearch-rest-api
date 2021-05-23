package app

import (
	"context"
	"fmt"
	"github.com/core-go/health"
	es "github.com/core-go/health/elasticsearch/v7"
	"github.com/core-go/log"
	"github.com/elastic/go-elasticsearch/v7"

	"go-service/internal/handlers"
	"go-service/internal/services"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	UserHandler   *handlers.UserHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	log.Initialize(root.Log)

	cfg := elasticsearch.Config{Addresses: []string{root.ElasticSearch.Url}}

	client, er1 := elasticsearch.NewClient(cfg)
	if er1 != nil {
		log.Error(ctx, "Cannot connect to elasticSearch. Error: "+er1.Error())
		return nil, er1
	}

	res, er2 := client.Info()
	if er2 != nil {
		log.Error(ctx, "Elastic server Error: " + er2.Error())
		return nil, er2
	}
	fmt.Println("Elastic server response: ", res)

	userService := services.NewUserService(client)
	userHandler := handlers.NewUserHandler(userService)

	elasticSearchChecker := es.NewHealthChecker(client)
	healthHandler := health.NewHandler(elasticSearchChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
