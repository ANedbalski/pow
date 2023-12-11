package pow

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	hashCahVersion = 1
	zeroByte       = 48
	randLimit      = 1000000000
)

var (
	ErrMaxIterationReached = errors.New("max iteration reached")
)

type TimeFunc func() time.Time

type HashCash struct {
	zeroBitCount int
	timeFunc     TimeFunc
}

func New(zeroCount int, timeFunc TimeFunc) *HashCash {
	return &HashCash{
		zeroBitCount: zeroCount,
		timeFunc:     timeFunc,
	}
}

func (hc HashCash) Validate(hash []byte) bool {
	if len(hash)*8 < hc.zeroBitCount {
		return false
	}

	f := hash[:hc.zeroBitCount/8+1]
	// make the irrelevant bits of the last byte zero
	f[len(f)-1] = f[len(f)-1] & byte(255<<(8-hc.zeroBitCount%8))

	for _, b := range f {
		if b != 0 {
			return false
		}
	}

	return true
}

// NewPuzzle returns a new puzzle implementing hashcash algorithm.
func (hc HashCash) NewPuzzle(sender string) Puzzle {
	return Puzzle{
		Version:    hashCahVersion,
		ZerosCount: hc.zeroBitCount,
		Date:       hc.timeFunc().Unix(),
		Sender:     sender,
		Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Intn(randLimit)))),
		Counter:    0,
	}
}

// BruteForce tries to find a solution for the given puzzle.
// If maxIterations is negative or zero, it will try to find a solution until it finds one.
// If maxIterations is positive, it will try to find a solution until it finds one or reaches maxIterations.
func (hc HashCash) BruteForce(puzzle Puzzle, maxIterations int) (Puzzle, error) {
	for !hc.Validate(puzzle.Hash()) {
		if maxIterations > 0 && puzzle.Counter > maxIterations {
			return Puzzle{}, ErrMaxIterationReached
		}
		puzzle.Counter++
	}
	return puzzle, nil
}
