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

-- name: MarkRequestReady :exec
UPDATE requests SET status = "Ready" WHERE id = ?;

-- name: MarkRequestSubmitted :exec
UPDATE requests SET status = "Submitted" WHERE id = ?;

-- name: MarkRequestInReview :exec
UPDATE requests SET status = "InReview", rpid = ? WHERE id = ?;

-- name: MarkRequestCanceled :exec
UPDATE requests SET status = "Canceled" WHERE id = ?;
