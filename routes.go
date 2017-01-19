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

var AllRoutes = map[string][]string{
	"/projects/": []string{"GET", "POST"},
	"/projects/{name}/": []string{"GET", "POST", "DELETE"},
	"/snapshot/": []string{"POST"},
}

func getFinalChar(s string) string {
	return string(s[len(s) - 1])
}

func WithTrailingSlash(rts Routes) Routes {
	var finalRoutes Routes

	for _, r := range rts {
		finalRoutes = append(finalRoutes, r)

		if getFinalChar(r.Pattern) != "/" {
			finalRoutes = append(finalRoutes, Route{
				Name: r.Name,
				Method: r.Method,
				Pattern: r.Pattern + "/",
				Handler: r.Handler,
			})
		}
	}
	return finalRoutes
}

/*
   Define Routes
       Handlers defined in handlers.go
*/
var routes = Routes{
	Route{
		"ListProjects",
		"GET",
		"/projects",
		ListProjects,
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
		"/projects",
		CreateProject,
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
		"/projects/{name}",
		DeleteProject,
	},
	Route{
		"SnapshotDB",
		"POST",
		"/snapshot",
		SnapshotDB,
	},
}
