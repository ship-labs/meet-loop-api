-- name: CreateGroup :one
INSERT INTO groups (name, description, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateGroupAdmin :one
INSERT INTO group_admins (group_id, member_id)
VALUES ($1, $2)
RETURNING *;

-- name: IsGroupAdmin :one
SELECT EXISTS(
    SELECT 1 FROM group_admins
    WHERE member_id = $1 AND group_id = $2
) AS is_admin;

-- name: GetUserGrops :many
WITH user_group_membership AS (
    SELECT id AS member_id, group_id FROM members WHERE members.user_id = $1
)
SELECT g.* FROM user_group_membership ugm
JOIN group_admins ga ON ga.group_id = ugm.group_id AND ga.member_id = ugm.member_id
JOIN groups g ON ga.group_id = g.id
LIMIT $2;
