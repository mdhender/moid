// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"github.com/mdhender/moid/internal/config"
	"github.com/mdhender/moid/internal/services"
)

type PaddleWebhook struct {
	Config *config.Config
	Paddle *services.Paddle
}
