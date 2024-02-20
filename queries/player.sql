-- name: CreatePlayer :execresult
INSERT INTO players (username, pw_hash) VALUES (?, ?);

-- name: UpdatePlayerPassword :execresult
UPDATE players SET pw_hash = ? WHERE id = ?;

-- name: GetPlayer :one
SELECT * FROM players WHERE id = ?;

-- name: GetPlayerByUsername :one
SELECT * FROM players WHERE username = ?;

-- name: GetPlayerUsername :one
SELECT (username) FROM players WHERE id = ?;

-- name: GetPlayerUsernameById :one
SELECT (username) FROM players WHERE id = ?;

-- name: SearchPlayersByUsername :many
SELECT * FROM players WHERE username LIKE ?;

-- name: CreatePlayerPermission :execresult
INSERT INTO player_permissions (name, pid, ipid) VALUES (?, ?, ?);

-- name: DeletePlayerPermission :exec
DELETE FROM player_permissions WHERE name = ? AND pid = ?;

-- name: ListPlayerPermissions :many
SELECT * FROM player_permissions WHERE pid = ?;

-- name: CreatePlayerPermissionIssuedChangeHistory :exec
INSERT INTO player_permission_change_history (name, pid, ipid) VALUES (?, ?, ?);

-- name: CreatePlayerPermissionRevokedChangeHistory :exec
INSERT INTO player_permission_change_history (name, pid, ipid, revoked) VALUES (?, ?, ?, true);

-- name: CreatePlayerSettings :exec
INSERT INTO player_settings (theme, pid) VALUES (?, ?);

-- name: GetPlayerSettings :one
SELECT * FROM player_settings WHERE pid = ?;

-- name: UpdatePlayerSettingsTheme :exec
UPDATE player_settings SET theme = ? WHERE pid = ?;
