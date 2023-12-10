-- name: CreateCharacterApplicationContent :execresult
INSERT INTO
  character_application_content 
  (gender, name, short_description, description, backstory, rid) 
VALUES 
  (?, ?, ?, ?, ?, ?);

-- name: GetCharacterApplicationContent :one
SELECT * FROM character_application_content WHERE id = ?;

-- name: GetCharacterApplicationContentForRequest :one
SELECT * FROM character_application_content WHERE rid = ?;

-- name: ListCharacterApplicationsForPlayer :many
SELECT
  sqlc.embed(character_application_content), sqlc.embed(requests)
FROM
  requests
JOIN
  character_application_content
ON
  requests.id = character_application_content.rid
WHERE
  requests.pid = ?
AND
  requests.type = "CharacterApplication"
AND
  requests.status != "Archived"
AND
  requests.status != "Canceled";

-- name: CountOpenCharacterApplicationsForPlayer :one
SELECT
  COUNT(*)
FROM
  requests
WHERE
  pid = ?
AND
  type = "CharacterApplication"
AND
  status != "Archived"
AND
  status != "Canceled";

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
  (gender, name, short_description, description, backstory, vid, rid)
VALUES
  (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateCharacterApplicationContent :exec
UPDATE 
  character_application_content
SET 
  gender = ?,
  name = ?,
  short_description = ?,
  description = ?,
  backstory = ?,
  vid = ?
WHERE
  rid = ?;

-- name: UpdateCharacterApplicationContentName :exec
UPDATE character_application_content SET name = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentGender :exec
UPDATE character_application_content SET gender = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentShortDescription :exec
UPDATE character_application_content SET short_description = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentDescription :exec
UPDATE character_application_content SET description = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentBackstory :exec
UPDATE character_application_content SET backstory = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentVersion :exec
UPDATE character_application_content SET vid = ? WHERE rid = ?;

