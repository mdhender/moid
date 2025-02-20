// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"fmt"
	"github.com/mdhender/moid/internal/config"
	"github.com/mdhender/moid/internal/server"
	"github.com/mdhender/semver"
	"log"
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

	srv, err := server.New(cfg, app.Routes(), context.Background())
	if err != nil {
		log.Fatal(err)
	}
	srv.Start()

	log.Printf("moid: shutting down after %v\n", time.Since(started))

	//log.Fatal(app.ServeCommand.Execute(app.Config.Server.Scheme, app.Config.Server.Host, app.Config.Server.Port, app.Routes(), context.Background()))
}
