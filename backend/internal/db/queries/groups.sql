-- name: ListGroups :many
SELECT id, group_name, page_sort FROM page_group ORDER BY id;

-- name: GetGroup :one
SELECT id, group_name, page_sort FROM page_group WHERE id = ?;

-- name: GetGroupForUpdate :one
SELECT id, group_name, page_sort FROM page_group WHERE id = ? FOR UPDATE;

-- name: CreateGroup :execlastid
INSERT INTO page_group (group_name, page_sort) VALUES (?, '[]');

-- name: RenameGroup :execrows
UPDATE page_group SET group_name = ? WHERE id = ?;

-- name: UpdateGroupPageSort :exec
UPDATE page_group SET page_sort = ? WHERE id = ?;

-- name: DeleteGroup :exec
DELETE FROM page_group WHERE id = ?;
