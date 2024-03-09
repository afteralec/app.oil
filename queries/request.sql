-- name: ListRequestsForPlayer :many
SELECT * FROM requests WHERE pid = ?;

-- name: GetRequest :one
SELECT * FROM requests WHERE id = ?;

-- name: CountOpenRequests :one
SELECT
  COUNT(*)
FROM
  requests
WHERE
  pid = ? AND status != "Archived" AND status != "Canceled";

-- name: CreateRequest :execresult
INSERT INTO requests (type, pid) VALUES (?, ?);

-- name: IncrementRequestVersion :exec
UPDATE requests SET vid = vid + 1 WHERE id = ?;

-- name: CreateHistoryForRequestStatusChange :exec
INSERT INTO 
  request_status_change_history
  (rid, vid, status, pid)
VALUES
  (?, (SELECT vid FROM requests WHERE requests.id = rid), (SELECT status FROM requests WHERE requests.id = rid), ?);

-- name: UpdateRequestStatus :exec
UPDATE requests SET status = ? WHERE id = ?;

-- name: UpdateRequestReviewer :exec
UPDATE requests SET rpid = ? WHERE id = ?;

-- name: CreateRequestComment :execresult
INSERT INTO
  request_comments (text, field, pid, rid, vid) 
VALUES
  (?, ?, ?, ?, (SELECT vid FROM requests WHERE requests.id = rid));

-- name: GetCommentWithAuthor :one
SELECT 
  sqlc.embed(players), sqlc.embed(request_comments)
FROM 
  request_comments 
JOIN
  players
ON
  request_comments.pid = players.id
WHERE 
  request_comments.id = ?;

-- name: ListCommentsForRequest :many
SELECT * FROM request_comments WHERE rid = ?;

-- name: ListCommentsForRequestWithAuthor :many
SELECT
  sqlc.embed(players), sqlc.embed(request_comments)
FROM
  request_comments
JOIN
  players
ON
  request_comments.pid = players.id
WHERE
  rid = ?;

-- name: CountUnresolvedComments :one
SELECT COUNT(*) FROM request_comments WHERE rid = ? AND resolved = false;

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
