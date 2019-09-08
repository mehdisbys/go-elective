package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Streamer interface {
	// 	LoadCountries()
	Start(w http.ResponseWriter, r *http.Request)
	SendCountries(w http.ResponseWriter, r *http.Request)
	SetProcessor(processor Processor)
}

type Processor func(duration time.Duration, input []string) <-chan []byte

type CountryServer struct {
	start     chan bool
	input     []string
	upgrader  *websocket.Upgrader
	interval  time.Duration
	processor Processor
}

func NewServer() Streamer {
	var upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	return &CountryServer{
		start:    make(chan bool),
		upgrader: upgrader,
		input:    []string{"a", "b", "c", "d"},
		interval: time.Millisecond * 1000,
	}
}

func (s *CountryServer) Start(w http.ResponseWriter, r *http.Request) {
	close(s.start)
	w.WriteHeader(http.StatusOK)
}

func (s *CountryServer) SetProcessor(processor Processor) {
	s.processor = processor
}

func (s *CountryServer) SendCountries(w http.ResponseWriter, r *http.Request) {
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
	for s := range s.processor(s.interval, s.input) {
		err := c.WriteMessage(websocket.TextMessage, s)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
