-- name: CreatePlayerPermissions :copyfrom
INSERT INTO player_permissions (permission, pid) VALUES (?, ?);

-- name: ListPlayerPermissions :many
SELECT * FROM player_permissions WHERE pid = ?;
