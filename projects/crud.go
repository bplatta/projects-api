package projects

import (
    "fmt"
    "gopkg.in/redis.v5"
)

type StatusError interface {
    IsCallerError() bool
}

type CRUDError struct {
    message string `json:"message"`
    operation string `json:"operation"`
    context string `json:"context"`
    callerError bool `json:"fatal"`
}

func (e *CRUDError) Error() string {
    return fmt.Sprintf("DB ERROR [%s]: %s. CONTEXT {%s}", e.operation, e.message, e.context)
}

func (e *CRUDError) IsCallerError() bool {
    return !e.callerError
}

type DBOptions struct {
    Address string
    Port string
    PoolSize int
    Password string
}

type DB struct {
    Options DBOptions
}

func (db *DB) GetClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr: db.Options.Address + ":" + db.Options.Port,
        Password: db.Options.Password,
        PoolSize: db.Options.PoolSize,
    })
}


/**
   Projects data CRUD
   Create - new project
   Read - project as JSON
   Update - project fields
   Delete - delete a project
 */

func (db *DB) Read(name string) (*Project, error) {
    data, err := db.GetClient().HGetAll(getProjectKey(name)).Result()

    if err == redis.Nil {
        return nil, &CRUDError{
            message: fmt.Sprintf("Project with name `%s` does not exist", name),
            operation: "READ",
            callerError: true,
        }
    } else if err != nil {
        return nil, &CRUDError{
            message: err.Error(),
            operation: "READ",
            callerError: false,
        }
    } else {
        var p Project
        p.SetFromData(data)
        return &p, nil
    }
}

func (db *DB) List() (*[]Project, error) {
    c := db.GetClient()
    var projectsList []Project
    projectSet, err := c.SMembers("projects").Result()

    if err != nil {
        return nil, &CRUDError{
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

            projectsList = append(projectsList, projS)

            if err != nil {
                return nil, &CRUDError{
                    message: "DB ERROR: Projects SET DNE",
                    operation: "LIST",
                    callerError: false,
                }
            }
        }
        return &projectsList, nil
    }
}

func (db *DB) Create(p Project) error {
    _, err := db.GetClient().HMSet(getProjectKey(p.Name), p.ToMap()).Result()

    if err != nil {
        return &CRUDError{
            message: err.Error(),
            operation: "LIST",
            callerError: false,
        }
    }
    return nil
}

func (db *DB) Update(name string, p Project) error {
    _, err := db.GetClient().HMSet(getProjectKey(p.Name), p.ToMap()).Result()

    if err != nil {
        return &CRUDError{
            message: err.Error(),
            operation: "LIST",
            callerError: false,
        }
    }
    return nil
}

func (db *DB) Delete(name string) error {
    return nil
}
