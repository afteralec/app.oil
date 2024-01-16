-- name: GetRoomImage :one
SELECT * FROM room_images WHERE id = ?;

-- name: GetRoomImageByName :one
SELECT * FROM room_images WHERE name = ?;

-- name: ListRoomImages :many
SELECT * FROM room_images;

-- name: CreateRoomImage :execresult
INSERT INTO room_images (name, title, description, size) VALUES (?, ?, ?, ?);

-- name: UpdateRoomImageName :exec
UPDATE room_images SET name = ? WHERE id = ?;

-- name: UpdateRoomImageTitle :exec
UPDATE room_images SET title = ? WHERE id = ?;

-- name: UpdateRoomImageDescription :exec
UPDATE room_images SET description = ? WHERE id = ?;

-- name: UpdateRoomImageSize :exec
UPDATE room_images SET size = ? WHERE id = ?;

-- name: GetRoom :one
SELECT * FROM rooms WHERE id = ?;

-- name: GetRoomByImageId :one
SELECT * FROM rooms WHERE riid = ?;

-- name: ListRooms :many
SELECT * FROM rooms;

-- name: ListRoomsWithImage :many
SELECT
  sqlc.embed(room_images), sqlc.embed(rooms)
FROM
  rooms
JOIN
  room_images
ON
  rooms.riid = room_images.id;

-- name: CreateRoom :execresult
INSERT INTO rooms (riid) VALUES (?);

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
