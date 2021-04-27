package main

import (
	"context"
	"fmt"
	"github.com/common-go/log"
	"github.com/gorilla/mux"
	config "go-service/configs"
	"go-service/internal/app"
	m "github.com/common-go/middleware"
	"net/http"
	"strconv"
)

func main() {
	/*var conf app.Root
	er1 := config.Load(&conf, "configs/config")
	if er1 != nil {
		panic(er1)
	}
	ctx := context.Background()

	app, er2 := app.NewApp(ctx, conf)
	if er2 != nil {
		panic(er2)
	}

	go health.Serve(conf.Server, app.HealthHandler)
	app.Consume(ctx, app.ConsumerHandler.Handle)*/
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

	var config app.Root
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
