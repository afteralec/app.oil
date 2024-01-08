-- name: GetPlayerSettings :one
SELECT * FROM player_settings WHERE pid = ?;

-- name: UpdatePlayerSettingsTheme :exec
UPDATE player_settings SET theme = ? WHERE pid = ?;
