// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/moid/internal/actions"
	"github.com/mdhender/moid/internal/domains"
	"github.com/mdhender/moid/internal/middlewares"
	"github.com/mdhender/moid/internal/responders"
	"github.com/mdhender/moid/internal/router"
	"html/template"
	"net/http"
	"path/filepath"
)

func (a *application) Routes() http.Handler {

	r := router.New(middlewares.Static(a.Config.Assets.Path))

	// public routes (no authentication required)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})
	r.Get("/home", a.Controllers.Home.Show)
	r.Get("/blogs", a.Controllers.Blogs.Show)
	r.Get("/reports", a.Controllers.Reports.Show)

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

	// Load templates
	tmpl := template.Must(template.ParseFiles(filepath.Join(a.Config.Views.Path, "user-row.gohtml")))

	// Dependency injection
	userRepo := &InMemoryUserRepo{data: make(map[string]domains.User)}
	userService := &domains.UserService{Repo: userRepo}
	createUserResponder := &responders.CreateUserResponder{Tmpl: tmpl}
	createUserAction := &actions.CreateUserAction{Service: userService, Responder: createUserResponder}

	// Register routes
	r.Post("/users", createUserAction.ServeHTTP)

	return r
}

// Mock repository implementation
type InMemoryUserRepo struct {
	data map[string]domains.User
}

func (repo *InMemoryUserRepo) Save(user domains.User) error {
	repo.data[user.Email] = user
	return nil
}
