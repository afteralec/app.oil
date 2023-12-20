-- name: CreateEmail :execresult
INSERT INTO emails (address, pid, verified) VALUES (?, ?, false);

-- name: MarkEmailVerified :execresult
UPDATE emails SET verified = true WHERE id = ?;

-- name: GetEmail :one
SELECT * FROM emails WHERE id = ?;

-- name: GetEmailByAddressForPlayer :one
SELECT * FROM emails WHERE address = ? AND pid = ?;

-- name: GetVerifiedEmailByAddress :one
SELECT * FROM emails WHERE address = ? AND verified = true;

-- name: ListEmails :many
SELECT * FROM emails WHERE pid = ?;

-- name: ListVerifiedEmails :many
SELECT * FROM emails WHERE pid = ? AND verified = true;

-- name: CountEmails :one
SELECT COUNT(*) FROM emails WHERE pid = ?;

-- name: DeleteEmail :execresult
DELETE FROM emails WHERE id = ?;
