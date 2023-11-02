-- name: CreateEmail :execresult
INSERT INTO emails (address, pid, verified) VALUES (?, ?, false);

-- name: MarkEmailVerified :execresult
UPDATE emails SET verified = true WHERE id = ?;

-- name: GetEmail :one
SELECT * FROM emails WHERE id = ?;

-- name: ListEmails :many
SELECT * FROM emails WHERE pid = ?;

-- name: CountEmails :one
SELECT COUNT(*) FROM emails WHERE pid = ?;

-- name: DeleteEmail :execresult
DELETE FROM emails WHERE id = ?;
