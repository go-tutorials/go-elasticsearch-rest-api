package app

import (
	"context"
	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	//PATCH  = "PATCH"
	DELETE = "DELETE"
)

func Route(r *mux.Router, ctx context.Context, config Config) error {
	app, err := NewApp(ctx, config)
	if err != nil {
		return err
	}

	r.HandleFunc("/health", app.HealthHandler.Check).Methods(GET)

	userPath := "/users"
	r.HandleFunc(userPath, app.UserHandler.All).Methods(GET)
	r.HandleFunc(userPath+"/{id}", app.UserHandler.Load).Methods(GET)
	r.HandleFunc(userPath, app.UserHandler.Insert).Methods(POST)
	r.HandleFunc(userPath+"/{id}", app.UserHandler.Update).Methods(PUT)
	r.HandleFunc(userPath+"/{id}", app.UserHandler.Delete).Methods(DELETE)

	return nil
}
