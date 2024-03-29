package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/mehdisbys/go-challenge/app/config"
	"github.com/mehdisbys/go-challenge/app/domain"
	"github.com/mehdisbys/go-challenge/app/server"
)

func main() {

	// Load config - it will exit if we are missing required values
	cfg, err := config.NewConfig()

	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}

	serve := server.NewServer(cfg.IntervalMs)

	serve.LoadCountries(cfg.CountriesFile)

	// The function that will send countries to the websocket
	serve.SetProcessor(domain.StreamValues)

	http.HandleFunc("/countries", serve.Start)
	http.HandleFunc("/events", serve.SendCountries)

	log.Printf("Listening on port %d", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil))
}
