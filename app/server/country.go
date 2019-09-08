package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Streamer interface {
	LoadCountries(filename string)
	SetProcessor(processor Processor)
	Start(w http.ResponseWriter, r *http.Request)
	SendCountries(w http.ResponseWriter, r *http.Request)
}

type Processor func(duration time.Duration, input []string) <-chan []byte

type CountryServer struct {
	start     chan bool
	input     []string
	upgrader  *websocket.Upgrader
	interval  time.Duration
	processor Processor
}

func NewServer(interval int) Streamer {
	var upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	return &CountryServer{
		start:    make(chan bool),
		upgrader: upgrader,
		interval: time.Millisecond * time.Duration(interval),
	}
}

func (s *CountryServer) LoadCountries(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	oneLine := string(content)
	s.input = strings.Split(oneLine, "\n")
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
