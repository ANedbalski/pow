package app

import (
	"context"
	"encoding/json"
	"errors"
	"pow/internal/book"
	"pow/internal/pow"
	"time"
)

type (
	Command string
	Message string
)

const (
	CmdGetPuzzle Command = "get_puzzle"
	CmdVerify    Command = "verify"
	CmdGetQuote  Command = "get_quote"

	MsgPuzzle   Message = "puzzle"
	MsgVerified Message = "verified"
	MsgQuote    Message = "quote"
	MsgError    Message = "error"
)

var (
	ErrNotVerified = errors.New("not verified")
	ErrInternal    = errors.New("internal error")
	ErrUnknownCmd  = errors.New("unknown command")
	ErrExpired     = errors.New("solution expired")
)

type POW interface {
	Validate(hash []byte) bool
	NewPuzzle(sender string) pow.Puzzle
}

type POWRepository interface {
	Add(key string, ttl time.Duration) error
	Exists(key string) (bool, error)
}

type Handler struct {
	pow         POW
	bookService *book.Service
	client      string
	verified    bool
	repo        POWRepository
	powTTL      time.Duration
}

func New(pow POW, bookService *book.Service, client string, repo POWRepository, powTTL time.Duration) *Handler {
	return &Handler{
		pow:         pow,
		bookService: bookService,
		client:      client,
		repo:        repo,
		powTTL:      powTTL,
	}
}

func (h *Handler) Dispatch(ctx context.Context, cmd Command, req string) (Message, string, error) {
	switch cmd {
	case CmdGetPuzzle:
		return h.getPuzzle(ctx)
	case CmdVerify:
		return h.verify(ctx, req)
	case CmdGetQuote:
		return h.getQuote(ctx, req)
	default:
		return MsgError, "", ErrUnknownCmd
	}
}

func (h *Handler) getPuzzle(_ context.Context) (Message, string, error) {
	puzzle := h.pow.NewPuzzle(h.client)
	msg, err := json.Marshal(puzzle)
	if err != nil {
		return MsgError, "", err
	}
	err = h.repo.Add(puzzle.Rand, h.powTTL)
	if err != nil {
		return MsgError, "", ErrInternal
	}

	return MsgPuzzle, string(msg), nil
}

func (h *Handler) verify(_ context.Context, payload string) (Message, string, error) {
	solution := pow.Puzzle{}
	err := json.Unmarshal([]byte(payload), &solution)
	if err != nil {
		return MsgError, "", ErrInternal
	}

	exists, err := h.repo.Exists(solution.Rand)
	if err != nil {
		return MsgError, "", ErrNotVerified
	}
	if !exists {
		return MsgError, "", ErrExpired
	}

	if !h.pow.Validate(solution.Hash()) {
		return MsgError, "", ErrNotVerified
	}

	h.verified = true

	return MsgVerified, "", nil
}

func (h *Handler) getQuote(ctx context.Context, _ string) (Message, string, error) {
	if !h.verified {
		return MsgError, "", ErrNotVerified
	}

	qoute, err := h.bookService.GetQuote(ctx)
	if err != nil {
		return MsgError, "", ErrInternal
	}

	return MsgQuote, qoute, nil
}
