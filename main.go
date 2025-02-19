// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"context"
	"github.com/mdhender/moid/internal/config"
	"github.com/mdhender/moid/internal/views"
	"github.com/mdhender/moid/ui"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile)
	cfg, err := config.Default(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	err = cfg.Load(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("config: %+v\n", cfg)

	app, err := newApplication(
		cfg,
		views.New(views.FS{
			FS:   ui.AssetsFS,
			Path: cfg.Views.AssetsPath,
		}, views.FS{
			FS:   ui.ViewsFS,
			Path: cfg.Views.ViewsPath,
		}),
	)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	log.Fatal(app.ServeCommand.Execute(app.Config.Server.Scheme, app.Config.Server.Host, app.Config.Server.Port, app.routes(), context.Background()))
}
