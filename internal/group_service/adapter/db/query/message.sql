-- name: AddMessage :one
INSERT INTO messages
(sender, group_id, text, status, created_at)
VALUES
($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetMessages :many
SELECT * FROM messages
WHERE group_id = $1
ORDER BY created_at
LIMIT $2
OFFSET $3;

-- name: UpdateMessageStatus :one
UPDATE messages
SET status = $2
WHERE id = $1
RETURNING *;