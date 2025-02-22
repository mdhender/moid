// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"github.com/mdhender/moid/internal/flash"
	"github.com/mdhender/moid/internal/sqlite"
	"github.com/mdhender/moid/internal/views"
	"log"
	"net/http"
)

type Reports struct {
	db   *sqlite.Store
	view *views.View
}

// NewReportsController creates a new instance of the Reports controller
func NewReportsController(db *sqlite.Store, view *views.View) (*Reports, error) {
	c := &Reports{
		db:   db,
		view: view,
	}
	// add any initialization logic here if needed
	return c, nil
}

func (c Reports) Show(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	store := flash.GetStore(r)

	// - Render the template
	c.view.Render(w, r, "reports.gohtml", struct {
		Error string
	}{
		Error: store.Get("error"),
	})
}
