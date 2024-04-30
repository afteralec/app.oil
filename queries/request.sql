-- name: ListRequestsForPlayer :many
SELECT * FROM requests WHERE pid = ?;

-- name: GetRequest :one
SELECT * FROM requests WHERE id = ?;

-- name: CreateRequest :execresult
INSERT INTO requests (type, status, pid) VALUES (?, ?, ?);

-- name: UpdateRequestStatus :exec
UPDATE requests SET status = ? WHERE id = ?;

-- name: UpdateRequestReviewer :exec
UPDATE requests SET rpid = ? WHERE id = ?;

-- name: CreateRequestField :exec
INSERT INTO request_fields (value, type, status, rid) VALUES (?, ?, ?, ?);

-- name: GetRequestField :one
SELECT * FROM request_fields WHERE id = ?;

-- name: GetRequestFieldByType :one
SELECT * FROM request_fields WHERE type = ? AND rid = ?;

-- name: GetRequestFieldByTypeWithChangeRequests :one
SELECT
  sqlc.embed(request_fields), sqlc.embed(open_request_change_requests), sqlc.embed(request_change_requests)
FROM
  request_fields
LEFT JOIN
  open_request_change_requests ON open_request_change_requests.rfid = request_fields.id
LEFT JOIN
  request_change_requests ON request_change_requests.rfid = request_fields.id
WHERE
  request_fields.type = ? AND request_fields.rid = ?;

-- name: ListRequestFieldsForRequest :many
SELECT * FROM request_fields WHERE rid = ?;

-- name: ListRequestFieldsForRequestWithChangeRequests :many
SELECT
  sqlc.embed(request_fields), sqlc.embed(open_request_change_requests), sqlc.embed(request_change_requests)
FROM
  request_fields
LEFT JOIN
  open_request_change_requests ON open_request_change_requests.rfid = request_fields.id
LEFT JOIN
  request_change_requests ON request_change_requests.rfid = request_fields.id
WHERE
  request_fields.rid = ?;

-- name: UpdateRequestFieldValue :exec
UPDATE request_fields SET value = ? WHERE id = ?;

-- name: UpdateRequestFieldValueByRequestAndType :exec
UPDATE request_fields SET value = ? WHERE type = ? AND rid = ?;

-- name: UpdateRequestFieldStatus :exec
UPDATE request_fields SET status = ? WHERE id = ?;

-- name: UpdateRequestFieldStatusByRequestAndType :exec
UPDATE request_fields SET status = ? WHERE type = ? AND rid = ?;

-- name: CreateOpenRequestChangeRequest :exec
INSERT INTO open_request_change_requests (value, text, rfid, pid) VALUES (?, ?, ?, ?);

-- name: GetOpenRequestChangeRequestForRequestField :one
SELECT * FROM open_request_change_requests WHERE rfid = ?;

-- name: CountOpenRequestChangeRequestsForRequest :one
SELECT
  COUNT(*)
FROM
  request_fields
JOIN
  open_request_change_requests ON open_request_change_requests.rfid = request_fields.id
WHERE
  request_fields.rid = ?;

-- name: GetOpenRequestChangeRequest :one
SELECT * FROM open_request_change_requests WHERE id = ?;

-- name: ListOpenRequestChangeRequestsByFieldID :many
SELECT * FROM open_request_change_requests WHERE rfid IN (sqlc.slice("rfids"));

-- name: DeleteOpenRequestChangeRequest :exec
DELETE FROM open_request_change_requests WHERE id = ?;

-- name: EditOpenRequestChangeRequest :exec
UPDATE open_request_change_requests SET text = ? WHERE id = ?;

-- name: CreateRequestChangeRequest :exec
INSERT INTO
  request_change_requests
SELECT * FROM
  open_request_change_requests
WHERE
  open_request_change_requests.id = ?;

-- name: DeleteRequestChangeRequest :exec
DELETE FROM request_change_requests WHERE id = ?;

-- name: CreatePastRequestChangeRequest :exec
INSERT INTO
  past_request_change_requests
SELECT * FROM
  request_change_requests
WHERE
  request_change_requests.id = ?;

-- name: ListRequestsByTypeAndStatus :many
SELECT * FROM requests WHERE type = ? AND status IN (sqlc.slice("statuses"));
