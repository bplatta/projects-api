package projects

import (
    "fmt"
    "gopkg.in/redis.v5"
    "errors"
)

const (
    PROJECTS = "projects"
)

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
        return nil, asCRUDError(
            errors.New(fmt.Sprintf("Project with name `%s` does not exist", name)),
            "READ", false)
    } else if err != nil {
        return nil, asCRUDError(err, "READ", true)
    } else {
        var p Project
        // Project by name does not exist
        if data["name"] == "" {
            return nil, nil
        }

        p.SetFromData(data)
        return &p, nil
    }
}

func (db *DB) List() (*[]Project, error) {
    c := db.GetClient()
    var projectsList []Project
    projectSet, err := c.SMembers(PROJECTS).Result()

    if err != nil {
        return nil, asCRUDError(errors.New("DB ERROR: Projects SET DNE"), "LIST", true)
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
                return nil, asCRUDError(errors.New("DB ERROR: Projects SET DNE"), "LIST", true)
            }
        }
        return &projectsList, nil
    }
}

func (db *DB) Create(p Project) error {
    key := getProjectKey(p.Name)
    c := db.GetClient()
    _, err := c.HMSet(key, p.ToMap()).Result()

    if err != nil {
        return asCRUDError(err, "LIST", true)
    }

    _, err = c.SAdd(PROJECTS, p.Name).Result()

    if err != nil {
        return asCRUDError(err, "LIST", true)
    }

    return nil
}

func (db *DB) Update(name string, p Project) error {
    _, err := db.GetClient().HMSet(getProjectKey(p.Name), p.ToMap()).Result()

    if err != nil {
        return asCRUDError(err, "Update", true)
    }
    return nil
}

func (db *DB) Delete(name string) error {
    c := db.GetClient()
    _, err := c.SRem(PROJECTS, name).Result()
    key := getProjectKey(name)
    _, err  = c.Del(key).Result()

    if err != nil {
        return asCRUDError(err, "LIST", true)
    }

    return nil
}
