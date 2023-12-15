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
INSERT INTO requests (type, pid, rpid) VALUES (?, ?, pid);

-- name: IncrementRequestVersion :exec
UPDATE requests SET vid = vid + 1 WHERE id = ?;

-- name: CreateHistoryForRequestStatus :exec
INSERT INTO 
  request_status_changes
  (rid, vid, status, pid)
VALUES
  (?, (SELECT vid FROM requests WHERE requests.rid = rid), (SELECT status FROM requests WHERE requests.rid = rid), ?);

-- name: MarkRequestReady :exec
UPDATE requests SET status = "Ready" WHERE id = ?;

-- name: MarkRequestSubmitted :exec
UPDATE requests SET status = "Submitted" WHERE id = ?;
