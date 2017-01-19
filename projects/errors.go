package projects

import (
    "fmt"
    "strings"
)

type StatusError interface {
    IsServerError() bool
}

type CRUDError struct {
    message string `json:"message"`
    operation string `json:"operation"`
    context string `json:"context"`
    ServerErr bool `json:"serverError"`
}

func (e *CRUDError) Error() string {
    return fmt.Sprintf("DB ERROR [%s]: %s. CONTEXT {%s}", e.operation, e.message, e.context)
}

func (e *CRUDError) IsServerError() bool {
    return !e.ServerErr
}

func asCRUDError(e error, op string, ServerErr bool) *CRUDError {
    return &CRUDError{
        message: e.Error(),
        operation: strings.ToUpper(op),
        ServerErr: ServerErr,
    }
}
