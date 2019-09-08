package main

import (
	"github.com/gorilla/websocket"
	"github.com/mehdisbys/go-challenge/app/domain"
	"log"
	"net/http"
	"time"
)

type Streamer interface {
	// 	LoadCountries()
	Start(w http.ResponseWriter, r *http.Request)
	SendCountries(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	start    chan bool
	input    []string
	upgrader *websocket.Upgrader
	interval time.Duration
}

func NewServer() Streamer {
	var upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	return &Server{
		start:    make(chan bool),
		upgrader: upgrader,
		input:    []string{"a", "b", "c", "d"},
		interval: time.Millisecond * 1000,
	}
}

func (s *Server) Start(w http.ResponseWriter, r *http.Request) {
	close(s.start)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) SendCountries(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer func() {
		err := c.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-s.start
	for s := range domain.StreamValues(s.interval, s.input) {
		err := c.WriteMessage(websocket.TextMessage, s)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func main() {
	serve := NewServer()
	http.HandleFunc("/countries", serve.Start)
	http.HandleFunc("/events", serve.SendCountries)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
