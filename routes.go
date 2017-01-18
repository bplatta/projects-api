package main

import (
	"github.com/bplatta/projects-api/projects"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Handler func(db *projects.DB, L Logger) http.Handler
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
	Route{
		"SnapshotDB",
		"POST",
		"/snapshot/",
		SnapshotDB,
	},
	Route{
		"SnapshotDB",
		"POST",
		"/snapshot",
		SnapshotDB,
	},
}
