package main

import (
	`github.com/mehdisbys/go-challenge/app/domain`
	`github.com/mehdisbys/go-challenge/app/server`
	"log"
	"net/http"
)

func main() {
	serve := server.NewServer()
	serve.SetProcessor(domain.StreamValues)
	http.HandleFunc("/countries", serve.Start)
	http.HandleFunc("/events", serve.SendCountries)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
