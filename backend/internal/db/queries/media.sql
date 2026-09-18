-- name: ListCarousels :many
SELECT id, image, url FROM carousel ORDER BY id;

-- name: GetCarousel :one
SELECT id, image, url FROM carousel WHERE id = ?;

-- name: CreateCarousel :execlastid
INSERT INTO carousel (image, url) VALUES (?, ?);

-- name: UpdateCarousel :execrows
UPDATE carousel SET image = ?, url = ? WHERE id = ?;

-- name: DeleteCarousel :execrows
DELETE FROM carousel WHERE id = ?;

-- name: FindCarouselsUsingText :many
SELECT id, image, url FROM carousel WHERE INSTR(image, sqlc.arg(needle)) > 0 ORDER BY id;

-- name: ListMarquees :many
SELECT id, text, color FROM marquee ORDER BY id;

-- name: GetMarquee :one
SELECT id, text, color FROM marquee WHERE id = ?;

-- name: CreateMarquee :execlastid
INSERT INTO marquee (text, color) VALUES (?, ?);

-- name: UpdateMarquee :execrows
UPDATE marquee SET text = ?, color = ? WHERE id = ?;

-- name: DeleteMarquee :execrows
DELETE FROM marquee WHERE id = ?;
