package main

import (
	"log"
	"net/http"
)

// main is the API entry point. Configures routes with Env settings.
// Serves data files as well as Projects API
func main() {
	config := GetConfigFromEnv()

	// construct router from routes.Routes configuration
	router := GetRouter(config)

	if config.HostData == true {
		FileDataRoute := "/data/"
		// serve raw data files
		router.PathPrefix(FileDataRoute).
			Handler(http.StripPrefix(FileDataRoute, http.FileServer(http.Dir(config.DataDir))))
	}

	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}
