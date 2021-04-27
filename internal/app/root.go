package app

import (
	"github.com/common-go/health"
	"github.com/common-go/kafka"
	"github.com/common-go/log"
	//"github.com/common-go/mongo"
	"github.com/common-go/mq"
	m "github.com/common-go/middleware"
)

type Root struct {
	Server      health.ServerConfig `mapstructure:"server"`
	Log         log.Config          `mapstructure:"log"`
	//Mongo       mongo.MongoConfig   `mapstructure:"mongo"`
	Retry       *mq.RetryConfig     `mapstructure:"retry"`
	Reader      ReaderConfig        `mapstructure:"reader"`
	KafkaWriter *kafka.WriterConfig `mapstructure:"kafka_writer"`
	MiddleWare m.LogConfig    `mapstructure:"middleware"`
}

type ReaderConfig struct {
	KafkaConsumer kafka.ReaderConfig `mapstructure:"kafka"`
	Config        mq.ConsumerConfig  `mapstructure:"retry"`
}

type ElasticClientConfig struct {
	UrlString string `mapstructure:"urlstring"`
}