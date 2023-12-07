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

-- name: UpdateCharacterApplicationContent :execresult
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
