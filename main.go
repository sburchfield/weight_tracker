package main

import (
	"net/http"
	"log"
	"time"
)

func main() {

	s := &http.Server{
		Addr:           ":"+envVars.appPort,
		Handler:        muxRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())

}
