package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	//"github.com/common-go/config"
	"github.com/common-go/log"
	m "github.com/common-go/middleware"
	"github.com/gorilla/mux"
	"go-service/configs"
	"go-service/internal/app"
)

func main() {
	var conf app.Root
	er1 := config.Load(&conf, "configs/config")
	if er1 != nil {
		panic(er1)
	}

	r := mux.NewRouter()

	log.Initialize(conf.Log)
	r.Use(m.BuildContext)
	logger := m.NewStructuredLogger()
	r.Use(m.Logger(conf.MiddleWare, log.InfoFields, logger))
	r.Use(m.Recover(log.ErrorMsg))

	var config app.ElasticClientConfig
	er2 := app.Route(r, context.Background(), config)
	if er2 != nil {
		panic(er2)
	}
	fmt.Println("Start server on port 5000")
	server := ""
	if conf.Server.Port > 0 {
		server = ":" + strconv.Itoa(5000)
	}
	http.ListenAndServe(server, r)
}
