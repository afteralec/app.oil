-- name: CreatePlayerPermission :execresult
INSERT INTO player_permissions (pid, permission) VALUES (?, ?);

-- name CreatePlayerPermissions :copyfrom
INSERT INTO player_permissions (pid, permission) VALUS (?, ?);

-- name: ListPlayerPermissions :many
SELECT * FROM player_permissions WHERE pid = ?;
