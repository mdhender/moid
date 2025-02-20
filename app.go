// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

import (
	"database/sql"
	"github.com/mdhender/moid/internal/actions"
	"github.com/mdhender/moid/internal/commands"
	"github.com/mdhender/moid/internal/config"
	"github.com/mdhender/moid/internal/controllers"
	"github.com/mdhender/moid/internal/encryption"
	"github.com/mdhender/moid/internal/ratelimiter"
	"github.com/mdhender/moid/internal/services"
	"github.com/mdhender/moid/internal/views"
	"path/filepath"
)

type application struct {
	Config *config.Config

	DB          *sql.DB
	Encrypter   *encryption.Encrypter
	RateLimiter *ratelimiter.Limiter
	Markdown    *services.Markdown
	Paddle      *services.Paddle

	Facades struct {
		Articles     *actions.ArticlesFacade
		Products     *actions.ProductsFacade
		Transactions *actions.TransactionsFacade
	}

	Controllers struct {
		Admin         *controllers.Admin
		Articles      *controllers.Articles
		Auth          *controllers.Auth
		Blogs         *controllers.Blogs
		Home          *controllers.Home
		Lqia          *controllers.Lqia
		PaddleWebhook *controllers.PaddleWebhook
		Ptg           *controllers.Ptg
		Purchases     *controllers.Purchases
		Reports       *controllers.Reports
	}

	Commands struct {
		Serve         *commands.Serve
		PaddleMigrate *commands.PaddleMigrate
	}

	Views *views.View
}

func newApplication(
	cfg *config.Config,
) (*application, error) {
	app := &application{
		Config: cfg,
	}

	// wire up the controllers for the application
	// should we be creating views for the controllers here?
	if blogsView, err := views.NewView("blogs.gohtml", filepath.Join(app.Config.Views.Path, "blogs.gohtml")); err != nil {
		return nil, err
	} else if app.Controllers.Blogs, err = controllers.NewBlogsController(app.DB, blogsView); err != nil {
		return nil, err
	}
	if homeView, err := views.NewView("home.gohtml", filepath.Join(app.Config.Views.Path, "home.gohtml")); err != nil {
		return nil, err
	} else if app.Controllers.Home, err = controllers.NewHomeController(app.DB, homeView); err != nil {
		return nil, err
	}
	if reportsView, err := views.NewView("reports.gohtml", filepath.Join(app.Config.Views.Path, "reports.gohtml")); err != nil {
		return nil, err
	} else if app.Controllers.Reports, err = controllers.NewReportsController(app.DB, reportsView); err != nil {
		return nil, err
	}

	return app, nil
}
