-- name: ListHelpSlugs :many
SELECT slug FROM help;

-- name: GetHelp :one
SELECT * FROM help WHERE slug = ?;

-- name: GetHelpRelated :many
SELECT * FROM help_related WHERE slug = ?;
