--  Copyright (c) 2025 Michael D Henderson. All rights reserved.
--

-- CreatePlayer creates a new player.
--
-- name: CreatePlayer :one
INSERT INTO players (name)
VALUES (:name)
RETURNING id;

-- CreateGame creates a new game.
--
-- name: CreateGame :one
INSERT INTO games (code, name, display_name)
VALUES (:code, :name, :display_name)
RETURNING id;

-- UpdateGameTurn increments the game turn number.
--
-- name: UpdateGameTurn :exec
UPDATE games
SET current_turn = :turn_number
WHERE id = :game_id;

-- GetCurrentGameTurn gets the current game turn.
--
-- name: GetCurrentGameTurn :one
SELECT current_turn
FROM games
WHERE id = :game_id;

-- CreateEmpire creates a new empire.
--
-- name: CreateEmpire :one
INSERT INTO empires (game_id, player_id)
VALUES (:game_id, :player_id)
RETURNING id;

-- CreateSystem creates a new system.
--
-- name: CreateSystem :one
INSERT INTO systems (game_id, x, y, z)
VALUES (:game_id, :x, :y, :z)
RETURNING id;

-- CreateStar creates a new star.
--
-- name: CreateStar :one
INSERT INTO stars (system_id, sequence)
VALUES (:system_id, :sequence)
RETURNING id;