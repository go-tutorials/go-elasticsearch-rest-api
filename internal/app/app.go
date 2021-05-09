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
	HealthHandler *health.HealthHandler
	UserHandler   *handlers.UserHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	log.Initialize(root.Log)

	cfg := elasticsearch.Config{
		Addresses: []string{root.ElasticSearch.Url},
		//Username: "<username>",
		//Password: "<password>",
	}

	client, er1 := elasticsearch.NewClient(cfg)

	if er1 != nil {
		log.Error(ctx, "Cannot connect to elasticSearch: Error: "+er1.Error())
		return nil, er1
	}

	res, err := client.Info()
	if err != nil {
		fmt.Println("Elastic server Error:", err)
	} else {
		fmt.Println("Elastic server response:", res)
	}

	userService := services.NewUserService(client)
	userHandler := handlers.NewUserHandler(userService)

	elasticSearchChecker := es.NewHealthChecker(client)
	healthHandler := health.NewHealthHandler(elasticSearchChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
