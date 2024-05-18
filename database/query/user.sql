-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: NewUser :one
INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $1 WHERE id = $2 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
