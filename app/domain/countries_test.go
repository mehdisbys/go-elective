package domain

import (
	"testing"
	"time"

	"github.com/go-test/deep"
)

func TestStreamValues(t *testing.T) {
	tests := []struct {
		name     string
		d        time.Duration
		input    []string
		expected [][]byte
	}{
		{
			name:     "Happy path",
			d:        time.Millisecond,
			input:    []string{"a", "b", "c"},
			expected: [][]byte{[]byte("a"), []byte("b"), []byte("c")},
		},
		{
			name:     "Empty input",
			d:        time.Millisecond,
			input:    []string{},
			expected: nil,
		},
		{
			name:     "Zero duration",
			d:        0,
			input:    []string{"a", "b", "c"},
			expected: [][]byte{[]byte("a"), []byte("b"), []byte("c")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := getValues(StreamValues(test.d, test.input))

			if diff := deep.Equal(test.expected, res); diff != nil {
				t.Error(diff)
			}

		})
	}
}

// getValues is a helper to append chan values into a slice
func getValues(c <-chan []byte) [][]byte {
	var r [][]byte

	for i := range c {
		r = append(r, i)
	}
	return r
}
