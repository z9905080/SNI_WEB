-- name: ListPageSummaries :many
SELECT id, page_group_id, page_name FROM page_content ORDER BY id;

-- name: ListPageIDsInGroup :many
SELECT id FROM page_content WHERE page_group_id = ? ORDER BY id;

-- name: CountPagesInGroup :one
SELECT COUNT(*) FROM page_content WHERE page_group_id = ?;

-- name: GetPage :one
SELECT id, page_group_id, page_name, html_context FROM page_content WHERE id = ?;

-- name: GetPageForUpdate :one
SELECT id, page_group_id, page_name, html_context FROM page_content WHERE id = ? FOR UPDATE;

-- name: CreatePage :execlastid
INSERT INTO page_content (page_group_id, page_name, html_context) VALUES (?, ?, ?);

-- name: UpdatePage :exec
UPDATE page_content SET page_group_id = ?, page_name = ?, html_context = ? WHERE id = ?;

-- name: DeletePage :exec
DELETE FROM page_content WHERE id = ?;

-- name: FindPagesUsingText :many
SELECT id, page_name FROM page_content WHERE INSTR(html_context, sqlc.arg(needle)) > 0 ORDER BY id;
