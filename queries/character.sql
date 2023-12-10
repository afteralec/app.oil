-- name: CreateCharacterApplicationContent :execresult
INSERT INTO
  character_application_content 
  (gender, name, sdesc, description, backstory, rid) 
VALUES 
  (?, ?, ?, ?, ?, ?);

-- name: GetCharacterApplicationContent :one
SELECT * FROM character_application_content WHERE id = ?;

-- name: GetCharacterApplicationContentForRequest :one
SELECT * FROM character_application_content WHERE rid = ?;

-- name: ListCharacterApplicationsForPlayer :many
SELECT
  sqlc.embed(requests), sqlc.embed(character_application_content)
FROM
  requests
JOIN
  character_application_content
ON
  requests.id = character_application_content.rid
WHERE
  requests.pid = ? AND requests.type = 'CharacterApplication';

-- name: ListCharacterApplicationContentForPlayer :many
SELECT
  *
FROM
  character_application_content 
WHERE
  rid
IN (SELECT id FROM requests WHERE pid = ?);

-- name: CreateCharacterApplicationContentHistory :execresult
INSERT INTO
  character_application_content_history
  (gender, name, sdesc, description, backstory, vid, rid)
VALUES
  (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateCharacterApplicationContent :exec
UPDATE 
  character_application_content
SET 
  gender = ?,
  name = ?,
  sdesc = ?,
  description = ?,
  backstory = ?,
  vid = ?
WHERE
  rid = ?;

-- name: UpdateCharacterApplicationContentName :exec
UPDATE character_application_content SET name = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentGender :exec
UPDATE character_application_content SET gender = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentSdesc :exec
UPDATE character_application_content SET sdesc = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentDescription :exec
UPDATE character_application_content SET description = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentBackstory :exec
UPDATE character_application_content SET backstory = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentVersion :exec
UPDATE character_application_content SET vid = ? WHERE rid = ?;

