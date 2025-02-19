// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package server provides the HTTP server for the application.
// The caller must provide the router for incoming HTTP requests
// and any needed middleware. This server implements a graceful
// shutdown of the server.
package server

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

// New returns a new instance of the server. Start the server
// with the Start method. Returns an error only if there is an
// invalid option passed in.
func New(cfg *config.Config, handler http.Handler, ctx context.Context) (*Server, error) {
	s := &Server{
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

	return s, nil
}

type Server struct {
	http.Server
	scheme string // must be http since we run behind a proxy
	host   string // should this be blank so that we're not bound to localhost?
	port   string
	ctx    context.Context
}

func (s *Server) BaseURL() string {
	return fmt.Sprintf("%s://%s", s.scheme, s.Addr)
}

// Start starts the server and blocks until the server is stopped.
// It implements a graceful shutdown of the server.
//
// TODO: the context should be used to cancel background tasks when shutting down the server.
func (s *Server) Start() {
	started := time.Now()

	// create a channel to listen for OS signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// start the server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("server: listening on %s\n", fmt.Sprintf("%s://%s", s.scheme, s.Addr))
		if err := http.ListenAndServe(s.Addr, s.Handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	ctxWithTimeout, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	// cancel any idle connections.
	log.Printf("server: canceling idle connections (%v)\n", time.Since(started))
	s.SetKeepAlivesEnabled(false)

	log.Printf("server: shutting down the server (%v)\n", time.Since(started))
	if err := s.Shutdown(ctxWithTimeout); err != nil {
		log.Fatalf("server: shutdown: %v\n", err)
	}

	log.Printf("server: ¡stopped gracefully! (%v)\n", time.Since(started))
}
