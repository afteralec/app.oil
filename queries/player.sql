-- name: CreatePlayer :execresult
INSERT INTO players (username, pw_hash) VALUES (?, ?);

-- name: GetPlayer :one
SELECT * FROM players WHERE id = ?;

-- name: GetPlayerByUsername :one
SELECT * FROM players WHERE username = ?;

-- name: GetPlayerPWHash :one
SELECT (pw_hash) FROM players WHERE id = ?;
