-- name: CreateRequestComment :execresult
INSERT INTO
  request_comments (text, field, pid, rid, vid) 
VALUES
  (?, ?, ?, ?, (SELECT vid FROM requests WHERE requests.id = rid));

-- name: ListCommentsForRequest :many
SELECT * FROM request_comments WHERE rid = ?;
