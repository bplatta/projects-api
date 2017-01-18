package main

import "net/http"

func LogRequestMiddleware(Log Logger, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        Log.LogRequest(r)
        next.ServeHTTP(w, r)
    })
}
