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

	cfg, err := config.NewConfig()

	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}

	serve := server.NewServer(cfg.IntervalMs)
	serve.LoadCountries(cfg.CountriesFile)
	serve.SetProcessor(domain.StreamValues)
	http.HandleFunc("/countries", serve.Start)
	http.HandleFunc("/events", serve.SendCountries)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil))
}
