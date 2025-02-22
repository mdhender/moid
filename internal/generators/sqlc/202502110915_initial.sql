--  Copyright (c) 2025 Michael D Henderson. All rights reserved.

-- foreign keys must be disabled to drop tables
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS meta_migrations;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS empires;
DROP TABLE IF EXISTS natural_resources;
DROP TABLE IF EXISTS orbits;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS stars;
DROP TABLE IF EXISTS systems;

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-- create the table for managing migrations
CREATE TABLE meta_migrations
(
    version    INTEGER  NOT NULL UNIQUE,
    comment    TEXT     NOT NULL,
    script     TEXT     NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- update the migrations table
INSERT INTO meta_migrations (version, comment, script)
VALUES (202502110915, 'initial migration', '202502110915_initial.sql');

CREATE TABLE players
(
    id INTEGER PRIMARY KEY AUTOINCREMENT
);

CREATE TABLE games
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    code         TEXT    NOT NULL UNIQUE,
    name         TEXT    NOT NULL UNIQUE,
    display_name TEXT    NOT NULL UNIQUE,
    current_turn INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE empires
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id   INTEGER NOT NULL REFERENCES games (id),
    player_id INTEGER NOT NULL REFERENCES players (id)
);

CREATE TABLE systems
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    game_id INTEGER NOT NULL REFERENCES games (id),
    x       INTEGER NOT NULL CHECK (x BETWEEN -15 AND 15),
    y       INTEGER NOT NULL CHECK (y BETWEEN -15 AND 15),
    z       INTEGER NOT NULL CHECK (z BETWEEN -15 AND 15)
);

CREATE TABLE stars
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    system_id INTEGER NOT NULL REFERENCES systems (id),
    sequence  INTEGER NOT NULL CHECK (sequence BETWEEN 1 AND 4),
    UNIQUE (system_id, sequence)
);

CREATE TABLE orbits
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    star_id      INTEGER NOT NULL REFERENCES stars (id),
    orbit        INTEGER NOT NULL CHECK (orbit BETWEEN 1 AND 10),
    kind         TEXT    NOT NULL CHECK (kind IN ('asteroid', 'empty', 'gas-giant', 'terrestrial')),
    habitability INTEGER NOT NULL CHECK (habitability BETWEEN 0 AND 25),
    UNIQUE (star_id, orbit)
);

CREATE TABLE natural_resources
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    orbit_id   INTEGER NOT NULL REFERENCES orbits (id),
    deposit_no INTEGER NOT NULL CHECK (deposit_no BETWEEN 1 AND 35),
    kind       TEXT    NOT NULL CHECK (kind IN ('fuel', 'gold', 'metallic', 'non-metallic')),
    quantity   INTEGER NOT NULL CHECK (quantity BETWEEN 0 AND 99000000),
    yield_pct  INTEGER NOT NULL CHECK (yield_pct BETWEEN 0 AND 100),
    UNIQUE (orbit_id, deposit_no)
);
