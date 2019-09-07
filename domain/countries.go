package domain

import (
	"time"
)

type SocketWriter interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
}

func StreamValues(output chan []byte, done chan bool, ticker *time.Ticker, strs []string) {
	defer ticker.Stop()

	for _, s := range strs {
		select {
		case <-ticker.C:
			output <- []byte(s)
		}
	}
	done <- true
}
