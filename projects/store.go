package projects

import (
    "gopkg.in/redis.v5"
)

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