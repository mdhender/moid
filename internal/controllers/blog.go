// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"database/sql"
	"github.com/mdhender/moid/internal/flash"
	"github.com/mdhender/moid/internal/views"
	"log"
	"net/http"
)

type Blog struct {
	db   *sql.DB
	view *views.View
}

// NewBlogController creates a new instance of the Blog controller
func NewBlogController(db *sql.DB, view *views.View) (*Blog, error) {
	c := &Blog{
		db:   db,
		view: view,
	}
	// add any initialization logic here if needed
	return c, nil
}

func (c Blog) Show(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	store := flash.GetStore(r)

	// - Render the template
	c.view.Render(w, r, "blog.gohtml", struct {
		Error string
	}{
		Error: store.Get("error"),
	})
}
