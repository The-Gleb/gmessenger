-- name: CreateSession :exec
INSERT INTO sessions (
    token, user_login, expiry
) VALUES (
    $1, $2, $3
);

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE token = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = $1;