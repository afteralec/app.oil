-- name: CreatePlayer :execresult
INSERT INTO players (username, pw_hash) VALUES (?, ?);

-- name: UpdatePlayerPassword :execresult
UPDATE players SET pw_hash = ? WHERE id = ?;

-- name: GetPlayer :one
SELECT * FROM players WHERE id = ?;

-- name: GetPlayerByUsername :one
SELECT * FROM players WHERE username = ?;

-- name: GetPlayerUsername :one
SELECT (username) FROM players WHERE username = ?;

-- name: GetPlayerUsernameById :one
SELECT (username) FROM players WHERE id = ?;

-- name: SearchPlayersByUsername :many
SELECT * FROM players WHERE username LIKE ?;
