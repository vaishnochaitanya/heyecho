-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :exec
INSERT INTO users (name, email, password, role, created_at) VALUES (?, ?, ?, ?, ?);

-- name: UpdateUser :one
UPDATE users SET name = ?, email = ?, password = ?, role = ? WHERE id = ? RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE email = ? RETURNING *;