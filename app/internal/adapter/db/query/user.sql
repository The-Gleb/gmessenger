-- name: CreateUser :one
INSERT INTO users
(username, login, password)
VALUES
($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE login = $1;

-- name: GetPassword :one
SELECT password FROM users
WHERE login = $1;

-- name: GetAllUsernames :many
SELECT username FROM users;

-- name: GetAllUsers :many
SELECT username, login FROM users;