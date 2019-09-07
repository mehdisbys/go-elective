package main

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	done := make(chan struct{})
	ticker := time.NewTicker(time.Second)
	for {
		writeToSocket(c, done, ticker)
	}
}

func writeToSocket(c *websocket.Conn, done chan struct{}, ticker *time.Ticker) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer ticker.Stop()
	i := 0
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(i)))
			if err != nil {
				log.Println("write:", err)
				return
			}
			i++
		}
	}
}

func main() {
	router := httprouter.New()
	router.GET("/countries", Countries)
	http.HandleFunc("/events", testWs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
