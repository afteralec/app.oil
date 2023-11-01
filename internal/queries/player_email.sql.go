// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: player_email.sql

package queries

import (
	"context"
	"database/sql"
)

const countPlayerEmails = `-- name: CountPlayerEmails :one
SELECT COUNT(*) FROM player_emails WHERE pid = ?
`

func (q *Queries) CountPlayerEmails(ctx context.Context, pid int64) (int64, error) {
	row := q.queryRow(ctx, q.countPlayerEmailsStmt, countPlayerEmails, pid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPlayerEmail = `-- name: CreatePlayerEmail :execresult
INSERT INTO player_emails (pid, email, verified) VALUES (?, ?, false)
`

type CreatePlayerEmailParams struct {
	Pid   int64
	Email string
}

func (q *Queries) CreatePlayerEmail(ctx context.Context, arg CreatePlayerEmailParams) (sql.Result, error) {
	return q.exec(ctx, q.createPlayerEmailStmt, createPlayerEmail, arg.Pid, arg.Email)
}

const deleteEmail = `-- name: DeleteEmail :execresult
DELETE FROM player_emails WHERE id = ?
`

func (q *Queries) DeleteEmail(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.deleteEmailStmt, deleteEmail, id)
}

const getEmail = `-- name: GetEmail :one
SELECT email, verified, pid, id FROM player_emails WHERE id = ?
`

func (q *Queries) GetEmail(ctx context.Context, id int64) (PlayerEmail, error) {
	row := q.queryRow(ctx, q.getEmailStmt, getEmail, id)
	var i PlayerEmail
	err := row.Scan(
		&i.Email,
		&i.Verified,
		&i.Pid,
		&i.ID,
	)
	return i, err
}

const listPlayerEmails = `-- name: ListPlayerEmails :many
SELECT email, verified, pid, id FROM player_emails WHERE pid = ?
`

func (q *Queries) ListPlayerEmails(ctx context.Context, pid int64) ([]PlayerEmail, error) {
	rows, err := q.query(ctx, q.listPlayerEmailsStmt, listPlayerEmails, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PlayerEmail
	for rows.Next() {
		var i PlayerEmail
		if err := rows.Scan(
			&i.Email,
			&i.Verified,
			&i.Pid,
			&i.ID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markEmailVerified = `-- name: MarkEmailVerified :execresult
UPDATE player_emails SET verified = true WHERE id = ?
`

func (q *Queries) MarkEmailVerified(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.markEmailVerifiedStmt, markEmailVerified, id)
}
