// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: player.sql

package queries

import (
	"context"
	"database/sql"
)

const createPlayer = `-- name: CreatePlayer :execresult
INSERT INTO players (username, role, pw_hash) VALUES (?, ?, ?)
`

type CreatePlayerParams struct {
	Username string
	Role     string
	PwHash   string
}

func (q *Queries) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (sql.Result, error) {
	return q.exec(ctx, q.createPlayerStmt, createPlayer, arg.Username, arg.Role, arg.PwHash)
}

const getPlayer = `-- name: GetPlayer :one
SELECT created_at, updated_at, pw_hash, username, role, id FROM players WHERE id = ?
`

func (q *Queries) GetPlayer(ctx context.Context, id int64) (Player, error) {
	row := q.queryRow(ctx, q.getPlayerStmt, getPlayer, id)
	var i Player
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PwHash,
		&i.Username,
		&i.Role,
		&i.ID,
	)
	return i, err
}

const getPlayerByUsername = `-- name: GetPlayerByUsername :one
SELECT created_at, updated_at, pw_hash, username, role, id FROM players WHERE username = ?
`

func (q *Queries) GetPlayerByUsername(ctx context.Context, username string) (Player, error) {
	row := q.queryRow(ctx, q.getPlayerByUsernameStmt, getPlayerByUsername, username)
	var i Player
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PwHash,
		&i.Username,
		&i.Role,
		&i.ID,
	)
	return i, err
}

const getPlayerUsername = `-- name: GetPlayerUsername :one
SELECT (username) FROM players WHERE username = ?
`

func (q *Queries) GetPlayerUsername(ctx context.Context, username string) (string, error) {
	row := q.queryRow(ctx, q.getPlayerUsernameStmt, getPlayerUsername, username)
	err := row.Scan(&username)
	return username, err
}

const getPlayerUsernameById = `-- name: GetPlayerUsernameById :one
SELECT (username) FROM players WHERE id = ?
`

func (q *Queries) GetPlayerUsernameById(ctx context.Context, id int64) (string, error) {
	row := q.queryRow(ctx, q.getPlayerUsernameByIdStmt, getPlayerUsernameById, id)
	var username string
	err := row.Scan(&username)
	return username, err
}

const updatePlayerPassword = `-- name: UpdatePlayerPassword :execresult
UPDATE players SET pw_hash = ? WHERE id = ?
`

type UpdatePlayerPasswordParams struct {
	PwHash string
	ID     int64
}

func (q *Queries) UpdatePlayerPassword(ctx context.Context, arg UpdatePlayerPasswordParams) (sql.Result, error) {
	return q.exec(ctx, q.updatePlayerPasswordStmt, updatePlayerPassword, arg.PwHash, arg.ID)
}
