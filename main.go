package gopiapi

import (
    "net/http"
    "log"
)

// main is the API entry point. Configures routes with Env settings.
// Serves data files as well as Projects API
func main() {

    FileDataRoute := "/data/"
    config := GetConfigFromEnv()

    // construct router from routes.Routes configuration
    router := GetRouter(config)

    // serve raw data files
    router.PathPrefix(FileDataRoute).
        Handler(http.StripPrefix(FileDataRoute, http.FileServer(http.Dir(config.DataDir))));

    log.Fatal(http.ListenAndServe(":" + config.Port, router))
}

