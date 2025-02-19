// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"database/sql"
	"github.com/mdhender/moid/internal/flash"
	"github.com/mdhender/moid/internal/ratelimiter"
	"github.com/mdhender/moid/internal/views"
	"log"
	"net/http"
)

type Home struct {
	db      *sql.DB
	limiter *ratelimiter.Limiter
	views   *views.View
}

// NewHomeController creates a new instance of the Home controller
func NewHomeController(db *sql.DB, limiter *ratelimiter.Limiter, views *views.View) (*Home, error) {
	c := &Home{
		db:    db,
		views: views,
	}
	// add any initialization logic here if needed
	return c, nil
}

func (c Home) Show(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	store := flash.GetStore(r)

	// TODO: Implement home page logic
	// - Fetch required data from database
	// - Process any markdown content
	// - Handle any necessary encryption/decryption
	// - Render the template

	c.views.Render(w, r, "index.gohtml", struct {
		Error string
	}{
		Error: store.Get("error"),
	})
}
