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

// Processor defines the function signature that will process the countries
type Processor func(duration time.Duration, input []string) <-chan []byte

// CountryServer implements Streamer. It synchronises websockets through a channel
type CountryServer struct {
	start     chan bool
	input     []string
	upgrader  *websocket.Upgrader
	interval  time.Duration
	processor Processor
}

// NewServer returns a Streamer
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

// LoadCountries loads the countries for later use
func (s *CountryServer) LoadCountries(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	oneLine := string(content)
	s.input = strings.Split(oneLine, "\n")
}

// Start gives the signal to all handlers to start writing the countries
func (s *CountryServer) Start(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	close(s.start)
	w.WriteHeader(http.StatusOK)
}

// SetProcessor sets the Processor that will return countries at different intervals
func (s *CountryServer) SetProcessor(processor Processor) {
	s.processor = processor
}

// SendCountries is a handler that will write the countries in a websocket.
// It calls the server's Processor function and writes the output.
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
