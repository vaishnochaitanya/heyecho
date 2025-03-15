-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserByEmailAndPassword :one
SELECT * FROM users WHERE email = ? AND password = ?;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :exec
INSERT INTO users (uuid, name, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateUser :one
UPDATE users SET name = ?, email = ?, password = ?, role = ?, updated_at = ? WHERE uuid = ? RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE uuid = ? RETURNING *;