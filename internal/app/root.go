package app

import (
	"github.com/common-go/log"
	m "github.com/common-go/middleware"
)

type Root struct {
	Server     ServerConfig   `mapstructure:"server"`
	DB         DatabaseConfig `mapstructure:"db"`
	Log        log.Config     `mapstructure:"log"`
	MiddleWare m.LogConfig    `mapstructure:"middleware"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver         string `mapstructure:"driver"`
	DataSourceName string `mapstructure:"data_source_name"`
}

type ElasticClientConfig struct {
	UrlString string `mapstructure:"urlstring"`
}
