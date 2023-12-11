package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"pow/config"
	"pow/internal/pow"
	"strings"
	"time"
)

func main() {
	cfg, err := config.NewClientConfig("./config", "client.yml")
	if err != nil {
		log.Fatalf("Error config file: %s \n", err)
	}

	conn, err := net.Dial("tcp", cfg.Server)
	panicOnErr(err)
	defer conn.Close()

	reader := bufio.NewReader(conn)

	//Request to get puzzle
	_, err = conn.Write([]byte("get_puzzle|{}\n"))
	panicOnErr(err)

	// Read puzzle from server
	getPuzzleMsg, err := reader.ReadString('\n')
	panicOnErr(err)
	_, j, err := parseMessage(getPuzzleMsg)
	panicOnErr(err)
	j = strings.Trim(j, "\n")
	puzzle := pow.Puzzle{}
	err = json.Unmarshal([]byte(j), &puzzle)
	panicOnErr(err)

	solver := pow.New(cfg.POWDifficulty, time.Now)
	solution, err := solver.BruteForce(puzzle, cfg.MaxIterations)
	panicOnErr(err)

	// Send puzzle solution to server
	msg, err := json.Marshal(solution)
	panicOnErr(err)
	_, err = conn.Write([]byte("verify|" + string(msg) + "\n"))
	panicOnErr(err)

	// Read quote from server
	verifyResp, err := reader.ReadString('\n')
	panicOnErr(err)
	cmd, _, err := parseMessage(verifyResp)
	panicOnErr(err)
	if cmd != "verified" {
		panic(errors.New("not verified"))
	}

	// Request quote from server
	_, err = conn.Write([]byte("get_quote|{}\n"))

	// Read quote from server
	getQuoteMsg, err := reader.ReadString('\n')
	_, quote, err := parseMessage(getQuoteMsg)
	panicOnErr(err)
	log.Printf("Quote: %s\n", quote)

}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func parseMessage(msg string) (string, string, error) {
	parts := strings.Split(msg, "|")
	if len(parts) != 2 {
		return "", "", errors.New("invalid message")
	}

	return parts[0], parts[1], nil
}
