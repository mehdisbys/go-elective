package main

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/mehdisbys/go-challenge/domain"
	"log"
	"net/http"
	"time"
)

func Countries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

func testWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	done := make(chan bool)
	output := make(chan []byte)
	ticker := time.NewTicker(time.Millisecond * 200)
	input := []string{"a", "b", "c", "d", "e"}
	for {
		go func(output chan []byte, done chan bool) {
			for {
				select {
				case <-done:
					c.WriteMessage(websocket.TextMessage, []byte("Bye!"))
					if err != nil {
						log.Println("write:", err)
						return
					}
					c.Close()
					return
				case s := <-output:
					err := c.WriteMessage(websocket.TextMessage, s)
					if err != nil {
						log.Println("write:", err)
						return
					}
				}
			}
		}(output, done)

		domain.StreamValues(output, done, ticker, input)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/countries", Countries)
	http.HandleFunc("/events", testWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
