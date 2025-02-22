// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/mdhender/moid/internal/config"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewServeHTTP returns a new instance of the serve command.
// Returns an error only if there is an invalid option passed in.
func NewServeHTTP(cfg *config.Config, handler http.Handler, ctx context.Context) (*ServeHTTP, error) {
	c := &ServeHTTP{
		scheme: "http",
		host:   cfg.Server.Host,
		port:   cfg.Server.Port,
		ctx:    ctx,
		Server: http.Server{
			Addr:           net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			Handler:        handler,
			ReadTimeout:    cfg.Server.ReadTimeout,
			WriteTimeout:   cfg.Server.WriteTimeout,
			IdleTimeout:    cfg.Server.IdleTimeout,
			MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
		},
	}
	return c, nil
}

func (c *ServeHTTP) BaseURL() string {
	return fmt.Sprintf("%s://%s", c.scheme, c.Addr)
}

// Execute starts the server and blocks until the server is stopped.
// It implements a graceful shutdown of the server.
//
// TODO: the context should be used to cancel background tasks when shutting down the server.
func (c *ServeHTTP) Execute() error {
	started := time.Now()

	// create a channel to listen for OS signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// start the server in a goroutine so that the server doesn't block.
	// note that WE will block and wait for signals to stop the server.
	go func() {
		log.Printf("server: listening on %s\n", fmt.Sprintf("%s://%s", c.scheme, c.Addr))
		if err := http.ListenAndServe(c.Addr, c.Handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %v\n", err)
		}
		log.Printf("server: shutdown\n")
	}()

	// server is running; block until we receive a signal.
	sig := <-stop
	log.Printf("server: signal %v: received (%v)\n", sig, time.Since(started))

	// graceful shutdown with a timeout.
	timeout := time.Second * 5
	log.Printf("server: timeout %v: creating context (%v)\n", timeout, time.Since(started))
	ctxWithTimeout, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	// cancel any idle connections.
	log.Printf("server: canceling idle connections (%v)\n", time.Since(started))
	c.SetKeepAlivesEnabled(false)

	log.Printf("server: shutting down the server (%v)\n", time.Since(started))
	if err := c.Shutdown(ctxWithTimeout); err != nil {
		return fmt.Errorf("server: shutdown: %w", err)
	}

	log.Printf("server: ¡stopped gracefully! (%v)\n", time.Since(started))
	return nil
}

type ServeHTTP struct {
	http.Server
	scheme string // must be http since we run behind a proxy
	host   string // should this be blank so that we're not bound to localhost?
	port   string
	ctx    context.Context
}
