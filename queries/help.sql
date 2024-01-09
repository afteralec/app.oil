-- name: ListHelpSlugs :many
SELECT slug FROM help;

-- name: ListHelpHeaders :many
SELECT slug, title, sub, category FROM help ORDER BY slug;

-- name: GetHelp :one
SELECT * FROM help WHERE slug = ?;

-- name: GetHelpRelated :many
SELECT * FROM help_related WHERE slug = ?;

-- name: SearchHelpByTitle :many
SELECT slug, title, sub FROM help WHERE slug LIKE ? OR title LIKE ?;

-- name: SearchHelpByContent :many
SELECT slug, title, sub FROM help WHERE raw LIKE ? OR sub LIKE ?;
