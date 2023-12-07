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

