package main

import (
	"github.com/bplatta/projects-api/projects"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(db *projects.DB) func(http.ResponseWriter, *http.Request)
}

type Routes []Route

/*
   Define Routes
       Handlers defined in handlers.go
*/
var routes = Routes{
	Route{
		"ListProjects",
		"GET",
		"/projects/",
		ListProjects,
	},
	Route{
		"ListProjects",
		"GET",
		"/projects",
		ListProjects,
	},
	Route{
		"ReadProject",
		"GET",
		"/projects/{name}/",
		ReadProject,
	},
	Route{
		"ReadProject",
		"GET",
		"/projects/{name}",
		ReadProject,
	},
	Route{
		"CreateProject",
		"POST",
		"/projects/",
		CreateProject,
	},
	Route{
		"CreateProject",
		"POST",
		"/projects",
		CreateProject,
	},
	Route{
		"UpdateProject",
		"POST",
		"/projects/{name}/",
		UpdateProject,
	},
	Route{
		"UpdateProject",
		"POST",
		"/projects/{name}",
		UpdateProject,
	},
	Route{
		"DeleteProject",
		"DELETE",
		"/projects/{name}/",
		DeleteProject,
	},
	Route{
		"DeleteProject",
		"DELETE",
		"/projects/{name}",
		DeleteProject,
	},
}

// GetRouter constructs a gorilla/mux router from
// routes var in the global scope. Requires a gopiapi.Config
// struct for configuration details. All Handlers
// are expected to accept a DB pointer and return a func of
// type http.HandlerFunc
func GetRouter(c Config) *mux.Router {

	router := mux.NewRouter()
	db := projects.DB{
		Options: projects.DBOptions{
			Address:  c.RedisHost,
			Port:     c.RedisPort,
			Password: c.RedisPassword,
			PoolSize: c.RedisPoolSize,
		},
	}

	for _, route := range routes {
		router.
			HandleFunc(route.Pattern, route.HandlerFunc(&db)).
			Name(route.Name).
			Methods(route.Method)
	}

	return router
}
