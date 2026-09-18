-- name: ListConfig :many
SELECT id, data_key, data_value FROM web_config ORDER BY id;

-- name: GetConfigForUpdate :one
SELECT id, data_value FROM web_config WHERE data_key = ? ORDER BY id LIMIT 1 FOR UPDATE;

-- name: InsertConfig :exec
INSERT INTO web_config (data_key, data_value) VALUES (?, ?);

-- name: UpdateConfig :exec
UPDATE web_config SET data_value = ? WHERE id = ?;
