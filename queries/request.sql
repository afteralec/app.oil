-- name: ListRequestsForPlayer :many
SELECT * FROM requests WHERE pid = ?;

-- name: GetRequest :one
SELECT * FROM requests WHERE id = ?;

-- name: ListCharacterApplicationsForPlayer :many
SELECT * FROM requests WHERE pid = ? AND type = 'CharacterApplication';

-- name: CreateRequest :execresult
INSERT INTO requests (type, pid) VALUES (?, ?);

-- name: CountOpenRequests :one
SELECT
  COUNT(*)
FROM
  requests
WHERE
  pid = ? AND status != "Archived" AND status != "Canceled";

-- TODO: Make these use the current vid of the application in the rid - compound query

-- name: AddCommentToRequest :execresult
INSERT INTO request_comments (text, pid, rid, vid) VALUES (?, ?, ?, ?);

-- name: AddCommentToRequestField :execresult
INSERT INTO request_comments (text, field, pid, rid, vid) VALUES (?, ?, ?, ?, ?);

-- TODO: Make this use the same field as the comment at the cid

-- name: AddReplyToComment :execresult
INSERT INTO request_comments (text, cid, pid, rid, vid) VALUES (?, ?, ?, ?, ?);

-- name: AddReplyToFieldComment :execresult
INSERT INTO request_comments (text, field, cid, pid, rid, vid) VALUES (?, ?, ?, ?, ?, ?);

-- name: ListCommentsForRequest :many
SELECT * FROM request_comments WHERE rid = ?;

-- name: GetRequestComment :one
SELECT * FROM request_comments WHERE id = ?;

-- name: ListRepliesToComment :many
SELECT * FROM request_comments WHERE cid = ?;

