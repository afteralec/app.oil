-- name: CreatePlayer :execresult
INSERT INTO players (username, role, pw_hash) VALUES (?, ?, ?);

-- name: GetPlayer :one
SELECT * FROM players WHERE id = ?;

-- name: GetRole :one
SELECT role FROM players WHERE id = ?;

-- name: GetPlayerByUsername :one
SELECT * FROM players WHERE username = ?;

-- name: GetPlayerUsername :one
SELECT (username) FROM players WHERE username = ?;

-- name: GetPlayerUsernameById :one
SELECT (username) FROM players WHERE id = ?;

-- name: GetPlayerPWHash :one
SELECT (pw_hash) FROM players WHERE id = ?;
