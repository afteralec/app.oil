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
