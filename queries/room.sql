-- name: GetRoom :one
SELECT * FROM rooms WHERE id = ?;

-- name: ListRooms :many
SELECT * FROM rooms;

-- name: ListRoomsByIDs :many
SELECT * FROM rooms WHERE id IN (sqlc.slice("ids"));

-- name: CreateRoom :execresult
INSERT INTO rooms (title, description, size) VALUES (?, ?, ?);

-- name: UpdateRoom :exec
UPDATE
  rooms
SET
  title = ?,
  description = ?,
  size = ?,
  unmodified = false
WHERE
  id = ?;

-- name: UpdateRoomTitle :exec
UPDATE rooms SET title = ?, unmodified = false WHERE id = ?;

-- name: UpdateRoomDescription :exec
UPDATE rooms SET description = ?, unmodified = false WHERE id = ?;

-- name: UpdateRoomSize :exec
UPDATE rooms SET size = ? WHERE id = ?;

-- name: UpdateRoomExitNorth :exec
UPDATE rooms SET north = ? WHERE id = ?;

-- name: UpdateRoomExitNortheast :exec
UPDATE rooms SET northeast = ? WHERE id = ?;

-- name: UpdateRoomExitEast :exec
UPDATE rooms SET east = ? WHERE id = ?;

-- name: UpdateRoomExitSoutheast :exec
UPDATE rooms SET southeast = ? WHERE id = ?;

-- name: UpdateRoomExitSouth :exec
UPDATE rooms SET south = ? WHERE id = ?;

-- name: UpdateRoomExitSouthwest :exec
UPDATE rooms SET southwest = ? WHERE id = ?;

-- name: UpdateRoomExitWest :exec
UPDATE rooms SET west = ? WHERE id = ?;

-- name: UpdateRoomExitNorthwest :exec
UPDATE rooms SET northwest = ? WHERE id = ?;
