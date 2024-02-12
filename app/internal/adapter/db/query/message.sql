-- name: CreateMessage :one
INSERT INTO messages
(sender, receiver, text, status, created_at)
VALUES
($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM messages
WHERE id = $1;

-- name: GetMessagesByUsers :many
SELECT * FROM messages
WHERE (sender = $1 AND receiver = $2)
OR (sender = $2 AND receiver = $1)
ORDER BY id
LIMIT $3 OFFSET $4;

-- name: GetLastMessage :one
SELECT * FROM messages
WHERE (sender = $1 AND receiver = $2)
OR (sender = $2 AND receiver = $1)
ORDER BY id
LIMIT $3 OFFSET $4;

-- name: UpdateMessageStatus :one
UPDATE messages
SET status = $1
WHERE id = $2
RETURNING *;

-- name: GetUnreadNumber :one
SELECT COUNT(*) FROM messages
WHERE (sender = $1 AND receiver = $2 AND status != 'READ');
