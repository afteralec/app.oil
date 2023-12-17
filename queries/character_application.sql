-- name: CreateCharacterApplicationContent :exec
INSERT INTO
  character_application_content 
  (gender, name, short_description, description, backstory, rid) 
VALUES 
  ("", "", "", "", "", ?);

-- name: GetCharacterApplicationContent :one
SELECT * FROM character_application_content WHERE id = ?;

-- name: GetCharacterApplication :one
SELECT
  sqlc.embed(character_application_content), sqlc.embed(requests)
FROM
  requests
JOIN
  character_application_content
ON
  character_application_content.rid = requests.id
WHERE
  requests.id = ?;

-- name: GetCharacterApplicationContentForRequest :one
SELECT * FROM character_application_content WHERE rid = ?;

-- name: ListCharacterApplicationsForPlayer :many
SELECT
  sqlc.embed(character_application_content), sqlc.embed(players), sqlc.embed(requests)
FROM
  requests
JOIN
  character_application_content
ON
  requests.id = character_application_content.rid
JOIN
  players
ON
  players.id = requests.pid
WHERE
  requests.pid = ?
AND
  requests.type = "CharacterApplication"
AND
  requests.status != "Archived"
AND
  requests.status != "Canceled";

-- name: ListOpenCharacterApplications :many
SELECT 
  sqlc.embed(character_application_content), sqlc.embed(players), sqlc.embed(requests)
FROM 
  requests
JOIN 
  character_application_content
ON 
  requests.id = character_application_content.rid
JOIN 
  players
ON 
  players.id = requests.pid
WHERE 
  requests.type = "CharacterApplication"
AND 
  requests.status = "Submitted"
OR 
  requests.status = "InReview"
OR 
  requests.status = "Reviewed";

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

-- name: CreateHistoryForCharacterApplication :exec
INSERT INTO
  character_application_content_history
  (gender, name, short_description, description, backstory, rid, vid)
SELECT 
  gender, name, short_description, description, backstory, rid, requests.vid
FROM
  character_application_content
JOIN
  requests
ON
  requests.id = character_application_content.rid
WHERE
  character_application_content.rid = ?;

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
