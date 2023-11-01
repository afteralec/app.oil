-- name: CreatePlayerEmail :execresult
INSERT INTO player_emails (pid, email, verified) VALUES (?, ?, false);

-- name: MarkEmailVerified :execresult
UPDATE player_emails SET verified = true WHERE id = ?;

-- name: GetEmail :one
SELECT * FROM player_emails WHERE id = ?;

-- name: ListPlayerEmails :many
SELECT * FROM player_emails WHERE pid = ?;

-- name: CountPlayerEmails :one
SELECT COUNT(*) FROM player_emails WHERE pid = ?;

-- name: DeleteEmail :execresult
DELETE FROM player_emails WHERE id = ?;
