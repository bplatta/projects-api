package gopiapi

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/bplatta/projects-api/projects"
)

type Route struct {
    Name string
    Method string
    Pattern string
    HandlerFunc func(db projects.DB) http.HandlerFunc
}

type Routes []Route

/*
    Define Routes
        Handlers defined in handlers.go
 */
var routes = Routes{
    Route{
        "ListModels",
        "GET",
        "/models",
        ListProjects,
    },
    Route{
        "ReadModel",
        "GET",
        "/models/{name}/",
        ReadProject,
    },
    Route{
        "CreateModel",
        "POST",
        "/models",
        CreateProject,
    },
    Route{
        "UpdateModel",
        "POST",
        "/models/{name}/",
        UpdateProject,
    },
    Route{
        "DeleteModel",
        "DELETE",
        "/models/{name}/",
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
        GetClient: projects.GetDBClient(&projects.DBOptions{
            Address: c.RedisHost,
            Port: c.RedisPort,
            Password: c.RedisPassword,
            PoolSize: c.RedisPoolSize,
        }),
    }

    for _, route := range routes {
        router.
            HandleFunc(route.Pattern, route.HandlerFunc(&db)).
            Name(route.Name).
            Methods(route.Method);
    }

    return &router
}
