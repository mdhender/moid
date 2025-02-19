// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Serve struct{}

func (c *Serve) Execute(scheme, host, port string, mux http.Handler, ctx context.Context) error {
	started := time.Now()
	log.Printf("server: schema %q: host %q: port %q\n", scheme, host, port)

	// create a channel to listen for OS signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	s := &http.Server{
		Addr:           net.JoinHostPort(host, port),
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    10 * time.Second,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// start the server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("listening on %s\n", fmt.Sprintf("%s://%s", scheme, s.Addr))
		if err := http.ListenAndServe(s.Addr, s.Handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %v\n", err)
		}
		log.Printf("server: shutdown\n")
	}()

	// server is running; block until we receive a signal.
	sig := <-stop
	log.Printf("signal: received %v (%v)\n", sig, time.Since(started))

	// graceful shutdown with a timeout.
	timeout := time.Second * 5
	log.Printf("creating context with %v timeout (%v)\n", timeout, time.Since(started))
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// cancel any idle connections.
	log.Printf("canceling idle connections (%v)\n", time.Since(started))
	s.SetKeepAlivesEnabled(false)

	log.Printf("sending signal to shut down the server (%v)\n", time.Since(started))
	if err := s.Shutdown(ctxWithTimeout); err != nil {
		log.Fatalf("server: shutdown: %v\n", err)
	}

	log.Printf("server stopped ¡gracefully! (%v)\n", time.Since(started))

	return nil
}
