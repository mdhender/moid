// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package db implements the data store via a Sqlite database.
//
// WARNING: depends on the generators/sqlc package to generate
// the database interfaces.
package db

import (
	"context"
	"database/sql"
	_ "embed"
)

type DB struct {
	path string
	db   *sql.DB
	ctx  context.Context
	q    *Queries
}
