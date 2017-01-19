package main

import (
	"encoding/json"
    "errors"
	"github.com/bplatta/projects-api/projects"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
    "fmt"
)

type ProjectListResp struct {
    Count int `json:"count"`
    Projects []projects.Project `json:"projects"`
}

/**
    API Routes
 */

// Lists the available endpoints and accepted methods
func ListRoutesRoot(routes map[string][]string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        asJSONResponse(w, routes, http.StatusOK)
    })
}

// ListProjects handles /projects/ route and retrieves the current
// projects, responding with JSON content-type. Panics if
// data cant be JSON encoded. All Handlers accept a DB struct and return
// handler with DB closure
func ListProjects(db *projects.DB, L Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if projectList, err := db.List(); err != nil {
			handleError(err, w, http.StatusServiceUnavailable) // 503
		} else {
            finalResp := map[string]interface{}{
                "count": len(*projectList),
                "projects": *projectList,
            }
            asJSONResponse(w, finalResp, http.StatusOK)
		}
	})
}

// ReadProject handles the GET request for a single Project. Panics
// if data or DB error can't be JSON encoded. All Handlers accept a DB struct and return
// handler with DB closure
func ReadProject(db *projects.DB, L Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := mux.Vars(r)["name"]
		if project, err := db.Read(n); err != nil {
			handleError(err, w, http.StatusServiceUnavailable) // 503
		} else {

            if project == nil {
                handleError(
                    errors.New(
                        fmt.Sprintf("Project with name `%s` does not exist", n)),
                    w, http.StatusNotFound)
            } else {
                asJSONResponse(w, *project, http.StatusOK)
            }
		}
	})
}

// CreateProject handles POST request for a new project.
// All Handlers accept a DB struct and return handler with DB closure. Required params: "name"
func CreateProject(db *projects.DB, L Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proj, err := convertBodyToProject(r)

        if err != nil {
			handleError(err, w, http.StatusUnprocessableEntity)
			return
		}

        if proj.Name == "" {
            handleError(errors.New("Missing required field: `name`"), w, http.StatusBadRequest)
            return
        }

		if err := db.Create(proj); err != nil {
			handleError(err, w, http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	})
}

// UpdateProject handles POST request to Project resources. Allow update
// of all Project fields. All Handlers accept a DB struct and return
// handler with DB closure
func UpdateProject(db *projects.DB, L Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proj, err := convertBodyToProject(r)

        if err != nil {
            handleError(err, w, http.StatusUnprocessableEntity)
			return
		}

		params := mux.Vars(r)
		n := params["name"]

		if err := db.Update(n, proj); err != nil {
			handleError(err, w, getStatusForError(err.(projects.StatusError)))
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	})
}

// DeleteProject handles DELETE request to Project resource. All Handlers
// accept a DB struct and return handler with DB closure
func DeleteProject(db *projects.DB, L Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		n := params["name"]

		if err := db.Delete(n); err != nil {
			handleError(err, w, getStatusForError(err.(projects.StatusError)))
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	})
}

func SnapshotDB(db *projects.DB, L Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
}

/**
    Route Handler helpers for managing errors bubbled up
 */

func asJSONResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(status)
    if encodeError := json.NewEncoder(w).Encode(data); encodeError != nil {
        panic(encodeError)
    }
}

// handleError writes error as JSON to http.ResponseWrite.
// If parsing error as JSON fails, panic occurs.
func handleError(err error, w http.ResponseWriter, status int) {
    errData := map[string]string{
        "error": err.Error(),
    }
    asJSONResponse(w, errData, status)
}

// getStatusForError maps a projects pkg error to status code int
// based on whether error is result of callee or Internal
func getStatusForError(err projects.StatusError) int {
    if err.IsServerError() {
        return http.StatusInternalServerError
    } else {
        return http.StatusBadRequest
    }
}

// convertBodyToProject accepts an http Request and converts the
// request Body to a Project struct. Panics on IO failure and returns
// error of the body is not parsable into a Project
func convertBodyToProject(r *http.Request) (projects.Project, error) {
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }

    var proj projects.Project
    if err := json.Unmarshal(body, &proj); err != nil {
        return proj, err
    }

    return proj, nil
}
