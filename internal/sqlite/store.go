// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package sqlite implements the data store via a Sqlite database.
//
// WARNING: depends on the generators/sqlc package to generate
// the database interfaces.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

type Store struct {
	path string
	db   *sql.DB
	ctx  context.Context
	q    *Queries
}

// Open opens the database at the given path.
// It returns an error if the path does not exist or is
// not a regular file.
func Open(path string, ctx context.Context) (*Store, error) {
	if abs, err := filepath.Abs(path); err != nil {
		return nil, err
	} else if sb, err := os.Stat(abs); err != nil {
		return nil, err
	} else if !sb.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file")
	} else {
		path = abs
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	log.Printf("store: %s: opened\n", path)
	return &Store{path: path, db: db, q: New(db), ctx: ctx}, nil
}

func (s *Store) Close() error {
	if s == nil && s.db == nil {
		return nil
	}
	defer func() {
		s.db = nil
		if s != nil {
			log.Printf("store: %s: closed\n", s.path)
		}
	}()
	return s.db.Close()
}
