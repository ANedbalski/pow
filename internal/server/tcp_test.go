package server

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net"
	"os"
	"pow/internal/app"
	"pow/internal/book"
	"pow/internal/pow"
	"pow/internal/repository"
	"strings"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	bookSrv := book.New(repository.NewInMemoBookRepo())
	powSrv := pow.New(20, time.Now)
	powRepo := repository.NewInMemoPOWRepo()
	srv := NewTCP("localhost:8080", logger, app.NewFactory(powSrv, bookSrv, powRepo, time.Minute))
	ctx, _ := context.WithCancel(context.Background())

	err := srv.Start(ctx)
	require.NoError(t, err)

	conn, err := net.Dial("tcp", "localhost:8080")
	require.NoError(t, err)
	reader := bufio.NewReader(conn)

	//Request to get puzzle
	_, err = conn.Write([]byte("get_puzzle|{}\n"))
	require.NoError(t, err)

	// Read puzzle from server
	getPuzzleMsg, err := reader.ReadString('\n')
	require.NoError(t, err)

	_, j, err := parseMessage(getPuzzleMsg)
	require.NoError(t, err)
	j = strings.Trim(j, "\n")
	puzzle := pow.Puzzle{}
	err = json.Unmarshal([]byte(j), &puzzle)
	require.NoError(t, err)

	solver := pow.New(20, time.Now)
	solution, err := solver.BruteForce(puzzle, 10000000)
	require.NoError(t, err)

	// Send puzzle solution to server
	msg, err := json.Marshal(solution)
	require.NoError(t, err)
	_, err = conn.Write([]byte("verify|" + string(msg) + "\n"))
	require.NoError(t, err)

	// Read verified from server
	verifyResp, err := reader.ReadString('\n')
	require.NoError(t, err)
	cmd, _, err := parseMessage(verifyResp)
	require.NoError(t, err)
	require.Equal(t, "verified", cmd)

	// Request quote from server
	_, err = conn.Write([]byte("get_quote|{}\n"))

	// Read quote from server
	getQuoteMsg, err := reader.ReadString('\n')
	_, _, err = parseMessage(getQuoteMsg)
	require.NoError(t, err)

	err = conn.Close()
	require.NoError(t, err)

	err = srv.Stop()
	require.NoError(t, err)
}

func parseMessage(msg string) (string, string, error) {
	parts := strings.Split(msg, "|")
	if len(parts) != 2 {
		return "", "", errors.New("invalid message")
	}

	return parts[0], parts[1], nil
}
