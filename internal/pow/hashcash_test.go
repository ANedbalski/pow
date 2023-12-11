package pow

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHashCash_BruteForce(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	thePuzzle := Puzzle{
		Version:    1,
		ZerosCount: 1,
		Date:       1701907521,
		Sender:     "sender",
		Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", 123460))),
		Counter:    0,
	}

	testCases := []struct {
		name          string
		zeroCount     int
		puzzle        Puzzle
		maxIterations int
		wantCount     int
		wantErr       error
	}{
		{
			name:          "zeroBitCount=6",
			zeroCount:     6,
			puzzle:        thePuzzle,
			maxIterations: -1,
			wantCount:     9,
			wantErr:       nil,
		},
		{
			name:          "zeroBitCount=12",
			zeroCount:     12,
			puzzle:        thePuzzle,
			maxIterations: -1,
			wantCount:     9981,
			wantErr:       nil,
		},
		{
			name:          "zeroBitCount=20",
			zeroCount:     20,
			puzzle:        thePuzzle,
			maxIterations: -1,
			wantCount:     864578,
			wantErr:       nil,
		},
		{
			name:          "limit iterations",
			zeroCount:     20,
			puzzle:        thePuzzle,
			maxIterations: 100,
			wantCount:     0,
			wantErr:       ErrMaxIterationReached,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.zeroCount, time.Now)

			got, err := h.BruteForce(tt.puzzle, tt.maxIterations)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantCount, got.Counter)
			fmt.Println(got.Hash())
		})
	}
}

func TestHashCash_Validate(t *testing.T) {
	testCases := []struct {
		name      string
		zeroCount int
		hash      []byte
		valid     bool
	}{
		{
			name:      "valid 1 zero bit",
			zeroCount: 1,
			hash:      []byte{0b01111111, 0b11111111, 0b11111111, 0b11111111},
			valid:     true,
		},
		{
			name:      "valid 8 zero bit",
			zeroCount: 8,
			hash:      []byte{0b00000000, 0b11111111, 0b11111111, 0b11111111},
			valid:     true,
		},
		{
			name:      "invalid 8 zero bit",
			zeroCount: 8,
			hash:      []byte{0b00000100, 0b11111111, 0b11111111, 0b11111111},
			valid:     false,
		},
		{
			name:      "valid 10 zero bit",
			zeroCount: 10,
			hash:      []byte{0b00000000, 0b00111111, 0b11111111, 0b11111111},
			valid:     true,
		},
		{
			name:      "invalid 10 zero bit",
			zeroCount: 10,
			hash:      []byte{0b00000100, 0b10111111, 0b11111111, 0b11111111},
			valid:     false,
		},
		{
			name:      "valid 20 zero bit",
			zeroCount: 20,
			hash:      []byte{0b00000000, 0b00000000, 0b00001111, 0b11111111},
			valid:     true,
		},
		{
			name:      "invalid 20 zero bit",
			zeroCount: 20,
			hash:      []byte{0b00000100, 0b10111111, 0b00001111, 0b11111111},
			valid:     false,
		},
		{
			name:      "invalid 20 zero bit",
			zeroCount: 20,
			hash:      []byte{0b00000000, 0b00000000, 0b00100000, 0b11111111},
			valid:     false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			h := New(tt.zeroCount, time.Now)
			assert.Equal(t, tt.valid, h.Validate(tt.hash), "should be valid", tt.hash)
		})
	}
}
