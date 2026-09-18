-- name: GetUserByAccount :one
SELECT id, account, pwd, identity, user_name FROM `user` WHERE account = ? ORDER BY id LIMIT 1;

-- name: GetUserByID :one
SELECT id, account, pwd, identity, user_name FROM `user` WHERE id = ?;

-- name: CreateUser :execlastid
INSERT INTO `user` (account, pwd, identity, user_name) VALUES (?, ?, 1, ?);

-- name: UpdateUserPassword :execrows
UPDATE `user` SET pwd = ? WHERE account = ?;

-- name: UpsertSession :exec
REPLACE INTO user_token (user_id, token, expire_time) VALUES (?, ?, ?);

-- name: GetSessionByToken :one
SELECT user_id, token, expire_time FROM user_token WHERE token = ? LIMIT 1;

-- name: ExtendSession :exec
UPDATE user_token SET expire_time = ? WHERE user_id = ?;

-- name: DeleteSession :exec
DELETE FROM user_token WHERE user_id = ?;
