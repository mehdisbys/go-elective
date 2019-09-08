package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// function p should not be called because Start() was not called
func TestSendCountriesNotCalled(t *testing.T) {
	var testNotcalled bool

	p := func(duration time.Duration, input []string) <-chan []byte {
		res := make(<-chan []byte)
		testNotcalled = true
		return res
	}

	s := NewServer(1)
	s.SetProcessor(p)

	// Create test server with the echo handler.
	testserver := httptest.NewServer(http.HandlerFunc(s.SendCountries))
	defer testserver.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(testserver.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	if testNotcalled {
		t.Error("function should not have been called")
	}

}

// function p should be called because Start() was called
func TestSendCountriesCalled(t *testing.T) {
	testCalled := true

	p := func(duration time.Duration, input []string) <-chan []byte {
		res := make(<-chan []byte)
		testCalled = false
		return res
	}

	s := NewServer(1)
	s.SetProcessor(p)

	// Create test server with the echo handler.
	testserver := httptest.NewServer(http.HandlerFunc(s.SendCountries))
	defer testserver.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(testserver.URL, "http")

	// Calling Start()
	s.Start(&httptest.ResponseRecorder{}, nil)

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	if testCalled {
		t.Error("function should have been called")
	}
}
