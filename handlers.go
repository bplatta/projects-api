package main

import (
	"encoding/json"
	"github.com/bplatta/projects-api/projects"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

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

// ListProjects handles /projects/ route and retrieves the current
// projects, responding with JSON content-type. Panics if
// data cant be JSON encoded. All Handlers accept a DB struct and return
// handler with DB closure
func ListProjects(db *projects.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if projectList, err := db.List(); err != nil {
			// 503
			w.WriteHeader(http.StatusServiceUnavailable)
			handleError(err, w)
		} else {
			if err := json.NewEncoder(w).Encode(projectList); err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}

// ReadProject handles the GET request for a single Project. Panics
// if data or DB error can't be JSON encoded. All Handlers accept a DB struct and return
// handler with DB closure
func ReadProject(db *projects.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		n := mux.Vars(r)["name"]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if project, err := db.Read(n); err != nil {
			// 503
			w.WriteHeader(http.StatusServiceUnavailable)
			handleError(err, w)
		} else {
			if err := json.NewEncoder(w).Encode(project); err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}

// CreateProject handles POST request for a new project.
// All Handlers accept a DB struct and return handler with DB closure
func CreateProject(db *projects.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		proj, err := convertBodyToProject(r)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			handleError(err, w)
			return
		}

		if err := db.Create(proj); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			handleError(err, w)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

// UpdateProject handles POST request to Project resources. Allow update
// of all Project fields. All Handlers accept a DB struct and return
// handler with DB closure
func UpdateProject(db *projects.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		proj, err := convertBodyToProject(r)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
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
	}
}

// DeleteProject handles DELETE request to Project resource. All Handlers
// accept a DB struct and return handler with DB closure
func DeleteProject(db *projects.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		n := params["name"]

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := db.Delete(n); err != nil {
			w.WriteHeader(getStatusForError(err.(projects.StatusError)))
			handleError(err, w)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
