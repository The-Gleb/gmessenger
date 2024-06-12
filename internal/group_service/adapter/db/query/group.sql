-- name: CreateGroup :one
INSERT INTO groups
(name, created_at)
VALUES
($1,$2)
RETURNING *;

-- name: GetGroups :many
SELECT * FROM
groups JOIN members
ON groups.id = members.group_id
WHERE member_login = $1;

-- name: Exists :one
SELECT CASE WHEN EXISTS (
    SELECT * FROM group
    WHERE id = $1
)
THEN TRUE
ELSE FALSE END;
