package app

import (
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
	sv "github.com/core-go/service"
)

type Root struct {
	Server        sv.ServerConfig     `mapstructure:"server"`
	ElasticSearch ElasticSearchConfig `mapstructure:"elastic_search"`
	Log           log.Config          `mapstructure:"log"`
	MiddleWare    mid.LogConfig       `mapstructure:"middleware"`
}

type ElasticSearchConfig struct {
	Url      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
