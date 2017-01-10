package projects

import (
    "gopkg.in/redis.v5"
)

type DBOptions struct {
    Address string
    Port string
    PoolSize int
    Password string
}

func GetDBClient(o *DBOptions) func() redis.Client {
    return func() redis.Client {
        return redis.NewClient(&redis.Options{
            Addr: o.Address + ":" + o.Port,
            Password: o.Password,
            PoolSize: o.PoolSize,
        })
    }
}

type ConnHandler struct {
    Client redis.Client
}

func (c ConnHandler) GetKey(k string) (string, error) {
    return "", nil
}

func (c ConnHandler) WriteToHashMap(key, field, value string) error {
    return nil
}

func (c ConnHandler) BackUpRDB(db int) (string, error) {
    return "", nil
}

func (c ConnHandler) LoadRDBFromDisk(fileName string) error {
    return nil
}