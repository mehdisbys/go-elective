package domain

import (
	"log"
	"time"
)

// StreamValues returns elements contained in `strs` at `d` intervals
func StreamValues(d time.Duration, strs []string) <-chan []byte {

	if d == 0 {
		// d cannot be zero, defaulting to lowest duration
		d = time.Nanosecond
	}

	ticker := time.NewTicker(d)
	rc := make(chan []byte)

	go func() {
		defer close(rc)
		defer ticker.Stop()

		for _, s := range strs {
			select {
			case <-ticker.C:
				log.Printf("Sending [%s] to websocket", s)
				rc <- []byte(s)
			}
		}
	}()
	return rc
}
