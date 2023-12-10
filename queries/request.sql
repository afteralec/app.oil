-- name: ListRequestsForPlayer :many
SELECT * FROM requests WHERE pid = ?;

-- name: GetRequest :one
SELECT * FROM requests WHERE id = ?;

-- name: CreateRequest :execresult
INSERT INTO requests (type, pid) VALUES (?, ?);

-- name: CountOpenRequests :one
SELECT
  COUNT(*)
FROM
  requests
WHERE
  pid = ? AND status != "Archived" AND status != "Canceled";

-- name: AddCommentToRequest :execresult
INSERT INTO request_comments (text, pid, rid, vid) VALUES (?, ?, ?, ?);

-- name: AddCommentToRequestField :execresult
INSERT INTO request_comments (text, field, pid, rid, vid) VALUES (?, ?, ?, ?, ?);

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

