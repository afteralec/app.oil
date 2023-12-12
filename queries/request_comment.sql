-- name: AddCommentToRequest :execresult
INSERT INTO
  request_comments (text, pid, rid, vid) 
VALUES
  (?, ?, ?, (SELECT vid FROM requests WHERE requests.rid = rid));

-- name: AddCommentToRequestField :execresult
INSERT INTO 
  request_comments (text, field, pid, rid, vid) 
VALUES 
  (?, ?, ?, ?, ?);

-- name: AddReplyToComment :execresult
INSERT INTO 
  request_comments (text, cid, pid, rid, vid) 
VALUES 
  (?, ?, ?, ?, (SELECT vid FROM requests WHERE requests.rid = rid));

-- name: AddReplyToFieldComment :execresult
INSERT INTO 
  request_comments (text, field, cid, pid, rid, vid) 
VALUES 
  (?, ?, ?, ?, ?, (SELECT vid FROM requests WHERE requests.rid = rid));

-- name: ListCommentsForRequest :many
SELECT * FROM request_comments WHERE rid = ?;

-- name: GetRequestComment :one
SELECT * FROM request_comments WHERE id = ?;

-- name: ListRepliesToComment :many
SELECT * FROM request_comments WHERE cid = ?;
