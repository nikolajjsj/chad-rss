-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password) VALUES (
  ?, ?
) RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET username = ?, password = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
