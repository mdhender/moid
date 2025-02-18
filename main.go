// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/moid/internal/config"
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

	_ = &application{}

	panic("!implemented")
}
