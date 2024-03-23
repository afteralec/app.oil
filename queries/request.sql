-- name: ListRequestsForPlayer :many
SELECT * FROM requests WHERE pid = ?;

-- name: GetRequest :one
SELECT * FROM requests WHERE id = ?;

-- name: CreateRequest :execresult
INSERT INTO requests (type, pid) VALUES (?, ?);

-- name: UpdateRequestStatus :exec
UPDATE requests SET status = ? WHERE id = ?;

-- name: UpdateRequestReviewer :exec
UPDATE requests SET rpid = ? WHERE id = ?;

-- name: CreateRequestField :exec
INSERT INTO request_fields (value, type, status, rid) VALUES (?, ?, ?, ?);

-- name: ListRequestFieldsForRequest :many
SELECT * FROM request_fields WHERE rid = ?;

-- name: UpdateRequestFieldValueByRequestAndType :exec
UPDATE request_fields SET value = ? WHERE type = ? AND rid = ?;

-- name: UpdateRequestFieldStatusByRequestAndType :exec
UPDATE request_fields SET status = ? WHERE type = ? AND rid = ?;

-- name: CreateRequestChangeRequest :exec
INSERT INTO request_change_requests (field, text, rid, pid) VALUES (?, ?, ?, ?);

-- name: GetRequestChangeRequest :one
SELECT * FROM request_change_requests WHERE id = ?;

-- name: DeleteRequestChangeRequest :exec
DELETE FROM request_change_requests WHERE id = ?;

-- name: EditRequestChangeRequest :exec
UPDATE request_change_requests SET text = ? WHERE id = ?;

-- name: GetCurrentRequestChangeRequestForRequestField :one
SELECT * FROM request_change_requests WHERE field = ? AND rid = ? AND old = false;

-- name: ListRequestChangeRequestsForRequest :many
SELECT * FROM request_change_requests WHERE rid = ? AND locked = ? AND old = ? ORDER BY updated_at;

-- name: ListRequestChangeRequestsForRequestField :many
SELECT * FROM request_change_requests WHERE field = ? AND rid = ? AND locked = ? AND old = ? ORDER BY updated_at;

-- name: ListCurrentRequestChangeRequestsForRequest :many
SELECT * FROM request_change_requests WHERE rid = ? AND old = false;

-- name: CountCurrentRequestChangeRequestForRequest :one
SELECT COUNT(*) FROM request_change_requests WHERE rid = ? AND old = false;

-- name: CountCurrentRequestChangeRequestForRequestField :one
SELECT COUNT(*) FROM request_change_requests WHERE field = ? AND rid = ? AND old = false;

-- name: LockRequestChangeRequestsForRequest :exec
UPDATE request_change_requests SET locked = true WHERE rid = ? AND locked = false;

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

-- name: CreateCharacterApplicationContentReview :exec
INSERT INTO
  character_application_content_review
  (gender, name, short_description, description, backstory, rid) 
VALUES 
  (?, ?, ?, ?, ?, ?);

-- name: GetCharacterApplicationContentReviewForRequest :one
SELECT * FROM character_application_content_review WHERE rid = ?;

-- name: UpdateCharacterApplicationContentReviewName :exec
UPDATE character_application_content_review SET name = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentReviewGender :exec
UPDATE character_application_content_review SET gender = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentReviewShortDescription :exec
UPDATE character_application_content_review SET short_description = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentReviewDescription :exec
UPDATE character_application_content_review SET description = ? WHERE rid = ?;

-- name: UpdateCharacterApplicationContentReviewBackstory :exec
UPDATE character_application_content_review SET backstory = ? WHERE rid = ?;
