package server

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"pow/internal/app"
	"strings"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, cmd app.Command, payload string) (app.Message, string, error)
}

type Connection struct {
	conn       net.Conn
	logger     *slog.Logger
	dispatcher Dispatcher
}

func NewConnection(conn net.Conn, d Dispatcher, l *slog.Logger) *Connection {
	return &Connection{
		conn:       conn,
		logger:     l,
		dispatcher: d,
	}
}

func (c *Connection) handle(ctx context.Context) {
	c.logger.Info("New connection: ", slog.String("remote_addr", c.conn.RemoteAddr().String()))
	defer c.conn.Close()

	reader := bufio.NewReader(c.conn)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Connection interrupted:")
			return
		default:
			req, err := reader.ReadString('\n')
			if errors.Is(err, io.EOF) {
				c.logger.Info("Connection closed by client")
				return
			}
			if err != nil {
				c.logger.Error("Read error", err)
				continue
			}
			resp, err := c.process(ctx, req)
			if err != nil {
				c.logger.Error("Process error", err)
				continue
			}

			err = c.respond(resp)
			if err != nil {
				c.logger.Error("Respond error", err)
			}
		}
	}
}

func (c *Connection) process(ctx context.Context, req string) (string, error) {
	cmd, payload, err := c.parse(req)
	if err != nil {
		return "", err
	}
	c.logger.Info("Receive request", slog.String("cmd", cmd), slog.String("payload", payload))

	msg, data, err := c.dispatcher.Dispatch(ctx, app.Command(cmd), payload)
	if err != nil {
		return string(msg) + "|" + err.Error(), nil
	}

	return string(msg) + "|" + data, err
}

func (c *Connection) parse(req string) (cmd string, r string, err error) {
	s := strings.Split(req, "|")
	if len(s) > 2 {
		return "", "", errors.New("invalid request")
	}
	if len(s) == 2 {
		r = s[1]
	}
	return s[0], strings.Trim(r, "\n"), nil
}

func (c *Connection) respond(msg string) error {
	msg = msg + "\n"
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		c.logger.Error("Write error", err)
	}

	return err
}
