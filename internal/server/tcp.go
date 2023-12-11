package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"pow/internal/app"
	"sync"
)

type HandlerFunc func(ctx context.Context, payload string) (string, error)

type DispatcherFactory interface {
	NewDispatcher(client string) *app.Handler
}

type TCP struct {
	Addr       string
	listener   net.Listener
	logger     *slog.Logger
	factory    DispatcherFactory
	stopSignal chan struct{}
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewTCP(addr string, logger *slog.Logger, df DispatcherFactory) *TCP {
	srv := &TCP{
		Addr:       addr,
		logger:     logger,
		factory:    df,
		stopSignal: make(chan struct{}),
	}

	return srv
}

func (t *TCP) Start(ctx context.Context) (err error) {
	if t.listener != nil {
		return errors.New("Server already running")
	}

	t.listener, err = net.Listen("tcp", t.Addr)
	if err != nil {
		return err
	}
	t.logger.Info("TCP server started", slog.String("addr", t.Addr))

	ctx, t.cancel = context.WithCancel(ctx)

	t.wg = sync.WaitGroup{}

	go t.listen(ctx)
	return nil
}

func (t *TCP) Stop() error {
	if t.listener == nil {
		t.logger.Warn("Stopping not running server")
		return nil
	}

	t.cancel()
	t.wg.Wait()

	err := t.listener.Close()
	return err
}

func (t *TCP) Running() {
	for {
		select {
		case <-t.stopSignal:
			return
		}
	}
}

func (t *TCP) listen(ctx context.Context) {
	defer func() {
		t.stopSignal <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := t.listener.Accept()
			if err != nil {
				t.logger.Error("TCP server accept error", err)
				return
			}

			dispatcher := t.factory.NewDispatcher(conn.RemoteAddr().String())
			connector := NewConnection(conn, dispatcher, t.logger)

			t.wg.Add(1)
			go func() {
				defer t.wg.Done()
				connector.handle(ctx)
			}()
		}
	}
}
