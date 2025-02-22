// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"fmt"
	"github.com/mdhender/moid/internal/config"
	"github.com/mdhender/semver"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	version = semver.Version{Major: 0, Minor: 0, Patch: 1}
)

func main() {
	// some hacks to get the version number
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("%s\n", version.String())
		os.Exit(0)
	}

	started := time.Now()

	log.SetFlags(log.Lshortfile)
	cfg, err := config.Default(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	err = cfg.Load(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	app, err := newApplication(
		cfg,
	)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	srv := &server{
		scheme: "http",
		host:   cfg.Server.Host,
		port:   cfg.Server.Port,
		ctx:    context.Background(),
		Server: http.Server{
			Addr:           net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
			Handler:        app.Routes(),
			ReadTimeout:    cfg.Server.ReadTimeout,
			WriteTimeout:   cfg.Server.WriteTimeout,
			IdleTimeout:    cfg.Server.IdleTimeout,
			MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
		},
	}
	if err = srv.start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("moid: shutting down after %v\n", time.Since(started))
}
