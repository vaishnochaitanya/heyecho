-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: ListUsers :many
SELECT * FROM users ORDER BY email;

-- name: CreateUser :exec
INSERT INTO users (email, password, role) VALUES (?, ?, ?);

-- name: UpdateUser :one
UPDATE users SET email = ?, password = ?, role = ? WHERE id = ? RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE email = ? RETURNING *;