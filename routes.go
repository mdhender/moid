// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/moid/internal/middlewares"
	"github.com/mdhender/moid/internal/router"
	"net/http"
)

func (a *application) Routes() http.Handler {

	r := router.New(middlewares.Static(a.Assets))

	// public routes (no authentication required)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})
	r.Get("/home", a.HomeController.Show)
	r.Get("/blog", a.BlogController.Show)

	//r := router.New(mid("zero"))
	//
	//r.Group(func(gr *router.Router) {
	//	gr.Use(mid("one"), mid("two"))
	//
	//	gr.Get("/foo", someHandler)
	//})
	//
	//r.Group(func(gr *router.Router) {
	//	gr.Use(mid("three"))
	//
	//	gr.Get("/bar", someHandler, mid("bar"))
	//	gr.Get("/baz", someHandler, mid("baz"))
	//})
	//
	//r.Post("/foobar", someHandler)
	//
	//// public routes (no authentication required)
	//r.Post("/login", handlers.LoginHandler)
	//r.Get("/refresh", handlers.RefreshTokenHandler)
	//r.Post("/logout", handlers.LogoutHandler)
	//
	//// protected routes (authentication required)
	//r.Group(func(gr *router.Router) {
	//	gr.Use(middleware.AuthMiddleware)
	//
	//	gr.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
	//		username := r.Context().Value("username").(string)
	//		_, _ = w.Write([]byte("Welcome, " + username + "\n"))
	//	})
	//})

	return r
}
