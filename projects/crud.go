package projects

import (
    "fmt"
    "gopkg.in/redis.v5"
)

type CRUDError struct {
    message string `json:"message"`
    operation string `json:"operation"`
    context string `json:"context"`
    callerError bool `json:"fatal"`
}

func (e CRUDError) Error() string {
    return fmt.Sprintf("DB ERROR [%s]: %s. CONTEXT {%s}", e.operation, e.message, e.context)
}

func (e CRUDError) IsCallerError() string {
    return !e.callerError
}

type DB struct {
    GetClient func() redis.Client
}

func (db *DB) Read(name string) (Project, CRUDError) {
    data, err := db.GetClient().HGetAll(getProjectKey(name)).Result()

    if err == redis.Nil {
        return nil, CRUDError{
            message: fmt.Sprintf("Project with name `%s` does not exist", name),
            operation: "READ",
            callerError: true,
        }
    } else if err != nil {
        return nil, CRUDError{
            message: err.Error(),
            operation: "READ",
            callerError: false,
        }
    } else {
        var p Project
        p.SetFromData(data)
        return p, nil
    }
}

func (db *DB) List() ([]Project, CRUDError) {
    c := db.GetClient()
    projectSet, err := c.SMembers("projects").Result()
    var projectsList []Project

    if err != nil {
        return nil, CRUDError{
            message: "DB ERROR: Projects SET DNE",
            operation: "LIST",
            callerError: false,
        }
    } else {
        for _, pName := range projectSet {
            proj, err := c.HGetAll(getProjectKey(pName)).Result()
            projS := Project{
                Name: pName,
                Model: proj["model"],
                DataLocation: proj["dataLocation"],
                URL: proj["url"],
            }

            append(projectsList, projS)

            if err != nil {
                return nil, CRUDError{
                    message: "DB ERROR: Projects SET DNE",
                    operation: "LIST",
                    callerError: false,
                }
            }
        }

        return projectsList
    }

}

func (db *DB) Create(p Project) CRUDError {
    _, err := db.GetClient().HMSet(getProjectKey(p.Name), p.ToMap()).Result()

    if err != nil {
        return nil, CRUDError{
            message: err.Error(),
            operation: "LIST",
            callerError: false,
        }
    }

    return nil
}

func (db *DB) Update(name string, p Project) CRUDError {
    _, err := db.GetClient().HMSet(getProjectKey(p.Name), p.ToMap()).Result()

    if err != nil {
        return nil, CRUDError{
            message: err.Error(),
            operation: "LIST",
            callerError: false,
        }
    }

    return nil
}

func (db *DB) Delete(name string) CRUDError {

}

