package app

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/go-playground/validator/v10"
	"go-service/internal/handlers"
	"go-service/internal/services"

	//"reflect"

	"github.com/common-go/health"
	"github.com/common-go/log"
	//"github.com/common-go/mongo"
	"github.com/common-go/mq"
	//v "github.com/common-go/validator"
	//"github.com/go-playground/validator/v10"
	//"github.com/sirupsen/logrus"
)

type ApplicationContext struct {
	HealthHandler   *health.HealthHandler
	UserHandler   	*handlers.UserHandler
	Consume         func(ctx context.Context, handle func(context.Context, *mq.Message, error) error)
	ConsumerHandler mq.ConsumerHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	log.Initialize(root.Log)

	//db, er1 := mongo.SetupMongo(ctx, root.Mongo)
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
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

	/*logError := log.ErrorMsg
	var logInfo func(context.Context, string)
	if logrus.IsLevelEnabled(logrus.InfoLevel) {
		logInfo = log.InfoMsg
	}*/

	/*consumer, er2 := kafka.NewReaderByConfig(root.Reader.KafkaConsumer, true)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new consumer. Error: "+er2.Error())
		return nil, er2
	}*/

	userService := services.NewEUserService(client)
	userHandler := handlers.NewUserHandler(userService)
	//elasticSearchChecker := health.NewelasticSearchHealthChecker(client)
	//checkers := []health.HealthChecker{elasticSearchChecker}
	//healthHandler := health.NewHealthHandler(checkers)
	/*userType := reflect.TypeOf(User{})
	writer := mongo.NewInserter(db, "users")
	validator := mq.NewValidator(userType, NewUserValidator().Validate)

	mongoChecker := mongo.NewHealthChecker(db)
	consumerChecker := kafka.NewKafkaHealthChecker(root.Reader.KafkaConsumer.Brokers, "kafka_consumer")
	var checkers []health.HealthChecker
	var consumerCaller mq.ConsumerHandler
	if root.KafkaWriter != nil {
		producer, er3 := kafka.NewWriterByConfig(*root.KafkaWriter)
		if er3 != nil {
			log.Error(ctx, "Cannot new a new producer. Error:"+er3.Error())
			return nil, er3
		}
		retryService := mq.NewMqRetryService(producer.Write, logError, logInfo)
		consumerCaller = mq.NewConsumerHandlerByConfig(root.Reader.Config, userType, writer.Write, retryService.Retry, validator.Validate, nil, logError, logInfo)
		producerChecker := kafka.NewKafkaHealthChecker(root.KafkaWriter.Brokers, "kafka_producer")
		checkers = []health.HealthChecker{mongoChecker, consumerChecker, producerChecker}
	} else {
		checkers = []health.HealthChecker{mongoChecker, consumerChecker}
		consumerCaller = mq.NewConsumerHandlerWithRetryConfig(userType, writer.Write, validator.Validate, root.Retry, true, logError, logInfo)
	}*/

	//handler := health.NewHealthHandler(checkers)
	return &ApplicationContext{
		//HealthHandler:   handler,
		UserHandler:   userHandler,
		//Consume:         consumer.Read,
		//ConsumerHandler: consumerCaller,
	}, nil
}

/*func NewUserValidator() v.Validator {
	validator := v.NewDefaultValidator()
	validator.CustomValidateList = append(validator.CustomValidateList, v.CustomValidate{Fn: CheckActive, Tag: "active"})
	return validator
}*/

func CheckActive(fl validator.FieldLevel) bool {
	return fl.Field().Bool()
}
