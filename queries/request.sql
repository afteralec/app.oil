-- name: ListRequestsForPlayer :many
SELECT * FROM requests WHERE pid = ?;

-- name: GetRequest :one
SELECT * FROM requests WHERE id = ?;

-- name: CreateRequest :execresult
INSERT INTO requests (type, status, pid) VALUES (?, ?, ?);

-- name: UpdateRequestStatus :exec
UPDATE requests SET status = ? WHERE id = ?;

-- name: UpdateRequestReviewer :exec
UPDATE requests SET rpid = ? WHERE id = ?;

-- name: CreateRequestField :exec
INSERT INTO request_fields (value, type, status, rid) VALUES (?, ?, ?, ?);

-- name: GetRequestField :one
SELECT * FROM request_fields WHERE id = ?;

-- name: GetRequestFieldByType :one
SELECT * FROM request_fields WHERE type = ? AND rid = ?;

-- name: GetRequestFieldByTypeWithChangeRequests :one
SELECT
  sqlc.embed(request_fields), sqlc.embed(open_request_change_requests), sqlc.embed(request_change_requests)
FROM
  request_fields
LEFT JOIN
  open_request_change_requests ON open_request_change_requests.rfid = request_fields.id
LEFT JOIN
  request_change_requests ON request_change_requests.rfid = request_fields.id
WHERE
  request_fields.type = ? AND request_fields.rid = ?;

-- name: ListRequestFieldsForRequest :many
SELECT * FROM request_fields WHERE rid = ?;

-- name: ListRequestFieldsForRequestWithChangeRequests :many
SELECT
  sqlc.embed(request_fields), sqlc.embed(open_request_change_requests), sqlc.embed(request_change_requests)
FROM
  request_fields
LEFT JOIN
  open_request_change_requests ON open_request_change_requests.rfid = request_fields.id
LEFT JOIN
  request_change_requests ON request_change_requests.rfid = request_fields.id
WHERE
  request_fields.rid = ?;

-- name: UpdateRequestFieldValue :exec
UPDATE request_fields SET value = ? WHERE id = ?;

-- name: UpdateRequestFieldValueByRequestAndType :exec
UPDATE request_fields SET value = ? WHERE type = ? AND rid = ?;

-- name: UpdateRequestFieldStatus :exec
UPDATE request_fields SET status = ? WHERE id = ?;

-- name: UpdateRequestFieldStatusByRequestAndType :exec
UPDATE request_fields SET status = ? WHERE type = ? AND rid = ?;

-- name: CreateOpenRequestChangeRequest :exec
INSERT INTO open_request_change_requests (value, text, rfid, pid) VALUES (?, ?, ?, ?);

-- name: GetOpenRequestChangeRequestForRequestField :one
SELECT * FROM open_request_change_requests WHERE rfid = ?;

-- name: CountOpenRequestChangeRequestsForRequestField :one
SELECT COUNT(*) FROM open_request_change_requests WHERE rfid = ?;

-- name: CountOpenRequestChangeRequestsForRequest :one
SELECT
  COUNT(*)
FROM
  request_fields
JOIN
  open_request_change_requests ON open_request_change_requests.rfid = request_fields.id
WHERE
  request_fields.rid = ?;

-- name: GetOpenRequestChangeRequest :one
SELECT * FROM open_request_change_requests WHERE id = ?;

-- name: DeleteOpenRequestChangeRequest :exec
DELETE FROM open_request_change_requests WHERE id = ?;

-- name: EditOpenRequestChangeRequest :exec
UPDATE open_request_change_requests SET text = ? WHERE id = ?;

-- name: CreateRequestChangeRequest :exec
INSERT INTO
  request_change_requests
SELECT * FROM
  open_request_change_requests
WHERE
  open_request_change_requests.id = ?;

-- name: DeleteRequestChangeRequest :exec
DELETE FROM request_change_requests WHERE id = ?;

-- name: CreatePastRequestChangeRequest :exec
INSERT INTO
  past_request_change_requests
SELECT * FROM
  request_change_requests
WHERE
  request_change_requests.id = ?;

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
