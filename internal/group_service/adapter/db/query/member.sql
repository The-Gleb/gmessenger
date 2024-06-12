-- name: AddMember :exec
INSERT INTO members
(group_id, member_login)
VALUES
($1,$2);

-- name: GetMembers :many
SELECT member_login FROM
groups JOIN members
ON groups.id = members.group_id
WHERE group_id = $1;

-- name: IsMember :one
SELECT CASE WHEN EXISTS (
    SELECT * FROM members
    WHERE member_login = $1 AND group_id = $2
)
THEN TRUE
ELSE FALSE END;
