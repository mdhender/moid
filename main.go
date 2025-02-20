// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"github.com/mdhender/moid/internal/config"
	"github.com/mdhender/moid/internal/server"
	"log"
	"os"
	"time"
)

func main() {
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

	srv, err := server.New(cfg, app.Routes(), context.Background())
	if err != nil {
		log.Fatal(err)
	}
	srv.Start()

	log.Printf("moid: shutting down after %v\n", time.Since(started))

	//log.Fatal(app.ServeCommand.Execute(app.Config.Server.Scheme, app.Config.Server.Host, app.Config.Server.Port, app.Routes(), context.Background()))
}
