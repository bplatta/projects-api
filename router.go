package main

import (
	"github.com/bplatta/projects-api/projects"
	"github.com/gorilla/mux"
)

// GetRouter constructs a gorilla/mux router from
// routes var in the global scope. Requires a gopiapi.Config
// struct for configuration details. All Handlers
// are expected to accept a DB pointer and return a func of
// type http.HandlerFunc
func GetRouter(c Config) *mux.Router {

	router := mux.NewRouter()
	logger := Logger{Level: c.LogLevel}
	db := projects.DB{
		Options: projects.DBOptions{
			Address:  c.RedisHost,
			Port:     c.RedisPort,
			Password: c.RedisPassword,
			PoolSize: c.RedisPoolSize,
		},
	}



	for _, route := range routes {
		// Pass DB reference and Logger to the Handler
		router.
			Handle(
				route.Pattern,
				LogRequestMiddleware(logger, route.Handler(&db, logger))).
			Name(route.Name).
			Methods(route.Method)
	}

	return router
}
