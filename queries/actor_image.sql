-- name: GetActorImage :one
SELECT * FROM actor_images WHERE id = ?;

-- name: ListActorImages :many
SELECT * FROM actor_images;

-- name: CreateActorImage :execresult
INSERT INTO actor_images (gender, name, short_description, description) VALUES (?, ?, ?, ?);

-- name: ListActorImageKeywords :many
SELECT * FROM actor_images_keywords WHERE aiid = ?;

-- name: CreateActorImageKeyword :execresult
INSERT INTO actor_images_keywords (aiid, keyword) VALUES (?, ?);

-- name: ListActorImageCan :many
SELECT * FROM actor_images_can WHERE aiid = ?;

-- name: CreateActorImageCan :execresult
INSERT INTO actor_images_can (aiid, can) VALUES (?, ?);

-- name: DeleteActorImageCan :exec
DELETE FROM actor_images_can WHERE id = ?;

-- name: ListActorImageCanBe :many
SELECT * FROM actor_images_can_be WHERE aiid = ?;

-- name: CreateActorImageCanBe :execresult
INSERT INTO actor_images_can_be (aiid, can_be) VALUES (?, ?);

-- name: DeleteActorImageCanBe :exec
DELETE FROM actor_images_can_be WHERE id = ?;

-- name: ListActorImagesHands :many
SELECT * FROM actor_images_hands WHERE aiid = ?;

-- name: CreateActorImageHand :execresult
INSERT INTO actor_images_hands (aiid, hand) VALUES (?, ?);

-- name: DeleteActorImageHand :exec
DELETE FROM actor_images_hands WHERE id = ?;

-- name: ListActorImagesPrimaryHands :many
SELECT * FROM actor_images_primary_hands WHERE aiid = ?;

-- name: CreateActorImagePrimaryHand :execresult
INSERT INTO actor_images_primary_hands (aiid, hand) VALUES (?, ?);

-- name: DeleteActorImagePrimaryHand :exec
DELETE FROM actor_images_primary_hands WHERE id = ?;

-- name: GetActorImageContainerProperties :one
SELECT * FROM actor_images_container_properties WHERE aiid = ?;

-- name: CreateActorImageContainerProperties :execresult
INSERT INTO actor_images_container_properties (aiid, is_container, is_surface_container, liquid_capacity) VALUES (?, ?, ?, ?);

-- name: DeleteActorImageContainerProperties :exec
DELETE FROM actor_images_container_properties WHERE id = ?;

-- name: GetActorImageFoodProperties :one
SELECT * FROM actor_images_food_properties WHERE aiid = ?;

-- name: CreateActorImageFoodProperties :execresult
INSERT INTO actor_images_food_properties (aiid, eats_into, sustenance) VALUES (?, ?, ?);

-- name: DeleteActorImageFoodProperties :exec
DELETE FROM actor_images_food_properties WHERE id = ?;

-- name: GetActorImageFurnitureProperties :one
SELECT * FROM actor_images_food_properties WHERE aiid = ?;

-- name: CreateActorImageFurnitureProperties :execresult
INSERT INTO actor_images_furniture_properties (aiid, seating) VALUES (?, ?);

-- name: DeleteActorImageFurnitureProperties :exec
DELETE FROM actor_images_furniture_properties WHERE id = ?;
