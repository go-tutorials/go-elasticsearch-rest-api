package app

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-playground/validator/v10"
	"go-service/internal/handlers"
	"go-service/internal/services"
	"github.com/common-go/health"
	"github.com/common-go/log"
	"github.com/common-go/mq"
)

type ApplicationContext struct {
	HealthHandler   *health.HealthHandler
	UserHandler   	*handlers.UserHandler
	Consume         func(ctx context.Context, handle func(context.Context, *mq.Message, error) error)
	ConsumerHandler mq.ConsumerHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	log.Initialize(root.Log)

	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
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

	userService := services.NewEUserService(client)
	userHandler := handlers.NewUserHandler(userService)

	return &ApplicationContext{
		UserHandler:   userHandler,
	}, nil
}

func CheckActive(fl validator.FieldLevel) bool {
	return fl.Field().Bool()
}
