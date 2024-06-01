-- name: GetActorImage :one
SELECT * FROM actor_images WHERE id = ?;

-- name: GetActorImageByName :one
SELECT * FROM actor_images WHERE name = ?;

-- name: ListActorImages :many
SELECT * FROM actor_images;

-- name: UpdateActorImageShortDescription :exec
UPDATE actor_images SET short_description = ? WHERE id = ?;

-- name: UpdateActorImageDescription :exec
UPDATE actor_images SET description = ? WHERE id = ?;

-- name: CreateActorImage :execresult
INSERT INTO actor_images (name, gender, short_description, description) VALUES (?, ?, ?, ?);

-- name: ListActorImageKeywords :many
SELECT * FROM actor_images_keywords WHERE aiid = ?;

-- name: CreateActorImageKeyword :execresult
INSERT INTO actor_images_keywords (keyword, aiid) VALUES (?, ?);

-- name: ListActorImageCan :many
SELECT * FROM actor_images_can WHERE aiid = ?;

-- name: CreateActorImageCan :execresult
INSERT INTO actor_images_can (can, aiid) VALUES (?, ?);

-- name: DeleteActorImageCan :exec
DELETE FROM actor_images_can WHERE id = ?;

-- name: ListActorImageCanBe :many
SELECT * FROM actor_images_can_be WHERE aiid = ?;

-- name: CreateActorImageCanBe :execresult
INSERT INTO actor_images_can_be (can_be, aiid) VALUES (?, ?);

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

-- name: CreateActorImageCharacterMetadata :exec
INSERT INTO actor_images_character_metadata (`key`, value, aiid) VALUES (?, ?, ?);

-- name: CreateActorImagePlayerProperties :execresult
INSERT INTO actor_images_player_properties (aiid, pid) VALUES (?, ?);

-- name: GetActorImagePlayerPropertiesForImage :one
SELECT * FROM actor_images_player_properties WHERE aiid = ?;

-- name: CountCurrentActorImagePlayerPropertiesForPlayer :one
SELECT COUNT(*) FROM actor_images_player_properties WHERE pid = ? AND current = true;

-- name: SetActorImagePlayerPropertiesCurrent :exec
UPDATE actor_images_player_properties SET current = ? WHERE id = ?;
