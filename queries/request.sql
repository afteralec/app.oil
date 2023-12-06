-- name: ListRequestsForPlayer :many
SELECT * FROM requests WHERE pid = ?;

-- name: GetRequest :one
SELECT * FROM requests WHERE id = ?;

-- name: ListCharacterApplicationsForPlayer :many
SELECT * FROM requests WHERE pid = ? AND type = 'CharacterApplication';

-- name: CreateRequest :execresult
INSERT INTO requests (pid) VALUES (?);

