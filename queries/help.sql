-- name: ListHelpSlugs :many
SELECT slug FROM help;

-- name: ListHelpHeaders :many
SELECT slug, title, sub, category FROM help ORDER BY slug;

-- name: GetHelp :one
SELECT * FROM help WHERE slug = ?;

-- name: GetHelpRelated :many
SELECT * FROM help_related WHERE slug = ?;

-- name: SearchHelpByTitle :many
SELECT slug, title, sub, category FROM help WHERE slug LIKE ? OR title LIKE ?;

-- name: SearchHelpByContent :many
SELECT slug, title, sub, category FROM help WHERE raw LIKE ? OR sub LIKE ?;

-- name: SearchHelpByCategory :many
SELECT slug, title, sub, category FROM help WHERE category LIKE ?;

-- name: SearchTags :many
SELECT slug, tag FROM help_tags WHERE tag LIKE ?;
