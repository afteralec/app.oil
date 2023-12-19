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
