package domain

import (
	`testing`
	`time`
)

func TestWriteToSocket(t *testing.T) {
	tests := []struct {
		name   string
		c      SocketWriter
		ticker *time.Ticker
		output chan []byte
		input  []string
	}{
		{

		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			StreamValues(test.output, nil, test.ticker, test.input)

		})
	}
}
