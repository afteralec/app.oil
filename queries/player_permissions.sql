-- name: CreatePlayerPermission :execresult
INSERT INTO player_permissions (pid, permission) VALUES (?, ?);

-- name: ListPlayerPermissions :many
SELECT * FROM player_permissions WHERE pid = ?;
