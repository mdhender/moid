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
	HomeController          *controllers.Home
	LqiaController          *controllers.Lqia
	PaddleWebhookController *controllers.PaddleWebhook
	PtgController           *controllers.Ptg
	PurchasesController     *controllers.Purchases

	ServeCommand         *commands.Serve
	PaddleMigrateCommand *commands.PaddleMigrate
}
