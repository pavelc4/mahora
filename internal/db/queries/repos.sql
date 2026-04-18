-- name: ListReposByUser :many
SELECT * FROM repos WHERE user_id = ? ORDER BY created_at DESC;

-- name: AddRepo :one
INSERT INTO repos (user_id, owner, name)
VALUES (?, ?, ?)
ON CONFLICT(user_id, owner, name) DO NOTHING
RETURNING *;

-- name: RemoveRepo :exec
DELETE FROM repos WHERE user_id = ? AND owner = ? AND name = ?;
