package main

import (
	"fmt"
	"log"
	"net/http"
)

const WEB_PORT = "8000"

type Config struct {
}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", WEB_PORT)

	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
