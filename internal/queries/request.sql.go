// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: request.sql

package queries

import (
	"context"
	"database/sql"
)

const countOpenRequests = `-- name: CountOpenRequests :one
SELECT
  COUNT(*)
FROM
  requests
WHERE
  pid = ? AND status != "Archived" AND status != "Canceled"
`

func (q *Queries) CountOpenRequests(ctx context.Context, pid int64) (int64, error) {
	row := q.queryRow(ctx, q.countOpenRequestsStmt, countOpenRequests, pid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createHistoryForRequestStatus = `-- name: CreateHistoryForRequestStatus :exec
INSERT INTO 
  request_status_changes
  (rid, vid, status, pid)
VALUES
  (?, (SELECT vid FROM requests WHERE requests.id = rid), (SELECT status FROM requests WHERE requests.id = rid), ?)
`

type CreateHistoryForRequestStatusParams struct {
	RID int64
	PID int64
}

func (q *Queries) CreateHistoryForRequestStatus(ctx context.Context, arg CreateHistoryForRequestStatusParams) error {
	_, err := q.exec(ctx, q.createHistoryForRequestStatusStmt, createHistoryForRequestStatus, arg.RID, arg.PID)
	return err
}

const createRequest = `-- name: CreateRequest :execresult
INSERT INTO requests (type, pid, rpid) VALUES (?, ?, pid)
`

type CreateRequestParams struct {
	Type string
	PID  int64
}

func (q *Queries) CreateRequest(ctx context.Context, arg CreateRequestParams) (sql.Result, error) {
	return q.exec(ctx, q.createRequestStmt, createRequest, arg.Type, arg.PID)
}

const getRequest = `-- name: GetRequest :one
SELECT created_at, updated_at, type, status, rpid, pid, id, vid, new FROM requests WHERE id = ?
`

func (q *Queries) GetRequest(ctx context.Context, id int64) (Request, error) {
	row := q.queryRow(ctx, q.getRequestStmt, getRequest, id)
	var i Request
	err := row.Scan(
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Type,
		&i.Status,
		&i.RPID,
		&i.PID,
		&i.ID,
		&i.VID,
		&i.New,
	)
	return i, err
}

const incrementRequestVersion = `-- name: IncrementRequestVersion :exec
UPDATE requests SET vid = vid + 1 WHERE id = ?
`

func (q *Queries) IncrementRequestVersion(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.incrementRequestVersionStmt, incrementRequestVersion, id)
	return err
}

const listRequestsForPlayer = `-- name: ListRequestsForPlayer :many
SELECT created_at, updated_at, type, status, rpid, pid, id, vid, new FROM requests WHERE pid = ?
`

func (q *Queries) ListRequestsForPlayer(ctx context.Context, pid int64) ([]Request, error) {
	rows, err := q.query(ctx, q.listRequestsForPlayerStmt, listRequestsForPlayer, pid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Request
	for rows.Next() {
		var i Request
		if err := rows.Scan(
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Type,
			&i.Status,
			&i.RPID,
			&i.PID,
			&i.ID,
			&i.VID,
			&i.New,
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

const markRequestCanceled = `-- name: MarkRequestCanceled :exec
UPDATE requests SET status = "Canceled" WHERE id = ?
`

func (q *Queries) MarkRequestCanceled(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.markRequestCanceledStmt, markRequestCanceled, id)
	return err
}

const markRequestReady = `-- name: MarkRequestReady :exec
UPDATE requests SET status = "Ready" WHERE id = ?
`

func (q *Queries) MarkRequestReady(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.markRequestReadyStmt, markRequestReady, id)
	return err
}

const markRequestSubmitted = `-- name: MarkRequestSubmitted :exec
UPDATE requests SET status = "Submitted" WHERE id = ?
`

func (q *Queries) MarkRequestSubmitted(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.markRequestSubmittedStmt, markRequestSubmitted, id)
	return err
}
