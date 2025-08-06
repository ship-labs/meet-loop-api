-- name: CreateGroupMember :one
INSERT INTO members (group_id, email, phone, name, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
