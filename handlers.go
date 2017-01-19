package main

import (
	"encoding/json"
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
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(routes)
    })
}

// ListProjects handles /projects/ route and retrieves the current
// projects, responding with JSON content-type. Panics if
// data cant be JSON encoded. All Handlers accept a DB struct and return
// handler with DB closure
func ListProjects(db *projects.DB, L Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if projectList, err := db.List(); err != nil {
			// 503
			w.WriteHeader(http.StatusServiceUnavailable)
			handleError(err, w)
		} else {
			w.WriteHeader(http.StatusOK)
            finalResp := map[string]interface{}{
                "count": len(*projectList),
                "projects": *projectList,
            }
			if err := json.NewEncoder(w).Encode(finalResp); err != nil {
				panic(err)
			}
		}
	})
}

// ReadProject handles the GET request for a single Project. Panics
// if data or DB error can't be JSON encoded. All Handlers accept a DB struct and return
// handler with DB closure
func ReadProject(db *projects.DB, L Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := mux.Vars(r)["name"]
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        e := json.NewEncoder(w)
		if project, err := db.Read(n); err != nil {
			// 503
			w.WriteHeader(http.StatusServiceUnavailable)
			handleError(err, w)
		} else {

            if project == nil {
                w.WriteHeader(http.StatusNotFound)
                e.Encode(
                    map[string]string{
                        "error": fmt.Sprintf("Project with name `%s` does not exist", n),
                    })
                return
            } else {
                w.WriteHeader(http.StatusOK)
                e.Encode(*project)
            }
		}
	})
}

// CreateProject handles POST request for a new project.
// All Handlers accept a DB struct and return handler with DB closure
func CreateProject(db *projects.DB, L Logger) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proj, err := convertBodyToProject(r)
        if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			handleError(err, w)
			return
		}

		if err := db.Create(proj); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			handleError(err, w)
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
			w.WriteHeader(http.StatusUnprocessableEntity)
            w.Header().Set("Content-Type", "application/json; charset=UTF-8")
            handleError(err, w)
			return
		}

		params := mux.Vars(r)
		n := params["name"]

		if err := db.Update(n, proj); err != nil {
			w.WriteHeader(getStatusForError(err.(projects.StatusError)))
			handleError(err, w)
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
			w.WriteHeader(getStatusForError(err.(projects.StatusError)))
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			handleError(err, w)
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

// handleError writes error as JSON to http.ResponseWrite.
// If parsing error as JSON fails, panic occurs.
func handleError(err error, w http.ResponseWriter) {
    if encodeError := json.NewEncoder(w).Encode(err); encodeError != nil {
        panic(encodeError)
    }
    return
}

// getStatusForError maps a projects pkg error to status code int
// based on whether error is result of callee or Internal
func getStatusForError(err projects.StatusError) int {
    if err.IsCallerError() {
        return http.StatusBadRequest
    } else {
        return http.StatusInternalServerError
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
