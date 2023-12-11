package pow

import (
	"crypto/sha1"
	"fmt"
)

type Puzzle struct {
	Version    int    `json:"version"`
	ZerosCount int    `json:"zeros_count"`
	Date       int64  `json:"date"`
	Sender     string `json:"sender"`
	Rand       string `json:"rand"`
	Counter    int    `json:"counter"`
}

func (p Puzzle) Hash() []byte {
	data := fmt.Sprintf("%d:%d:%d:%s:%s:%d", p.Version, p.ZerosCount, p.Date, p.Sender, p.Rand, p.Counter)
	h := sha1.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}
