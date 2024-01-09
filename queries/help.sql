-- name: ListHelpSlugs :many
SELECT slug FROM help;

-- name: ListHelpTitleAndSub :many
SELECT slug, title, sub FROM help ORDER BY slug;

-- name: GetHelp :one
SELECT * FROM help WHERE slug = ?;

-- name: GetHelpRelated :many
SELECT * FROM help_related WHERE slug = ?;
