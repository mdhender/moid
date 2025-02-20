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
)

type application struct {
	Config *config.Config

	DB          *sql.DB
	Encrypter   *encryption.Encrypter
	RateLimiter *ratelimiter.Limiter
	Markdown    *services.Markdown
	Paddle      *services.Paddle

	ArticlesFacade     *actions.ArticlesFacade
	ProductsFacade     *actions.ProductsFacade
	TransactionsFacade *actions.TransactionsFacade

	AdminController         *controllers.Admin
	ArticlesController      *controllers.Articles
	AuthController          *controllers.Auth
	BlogController          *controllers.Blog
	HomeController          *controllers.Home
	LqiaController          *controllers.Lqia
	PaddleWebhookController *controllers.PaddleWebhook
	PtgController           *controllers.Ptg
	PurchasesController     *controllers.Purchases

	ServeCommand         *commands.Serve
	PaddleMigrateCommand *commands.PaddleMigrate

	Views  *views.View
	Assets string
}

func newApplication(
	cfg *config.Config,
) (*application, error) {
	app := &application{
		Config: cfg,
		Assets: cfg.Assets.Path,
	}

	// wire up the controllers for the application
	// should we be creating views for the controllers here?
	if blogView, err := views.NewView("blog.gohtml", "../ui/views/blog.gohtml"); err != nil {
		return nil, err
	} else if app.BlogController, err = controllers.NewBlogController(app.DB, blogView); err != nil {
		return nil, err
	}
	if homeView, err := views.NewView("home.gohtml", "../ui/views/home.gohtml"); err != nil {
		return nil, err
	} else if app.HomeController, err = controllers.NewHomeController(app.DB, homeView); err != nil {
		return nil, err
	}

	return app, nil
}
