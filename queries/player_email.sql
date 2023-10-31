-- name: CreatePlayerEmail :execresult
INSERT INTO player_emails (pid, email, verified) VALUES (?, ?, false);

-- name: ListPlayerEmails :many
SELECT * FROM player_emails WHERE pid = ?;

-- name: CountPlayerEmails :one
SELECT COUNT(*) FROM player_emails WHERE pid = ?;
