-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
		gen_random_uuid(),
		NOW(),
		NOW(),
		$1
)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: CreateUserWithPassword :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
)
RETURNING *;

-- name: SetUserPassword :exec
UPDATE users
SET hashed_password = $2
WHERE users.id = $1;

-- name: GetUserWithEmail :one
Select * from users
WHERE email = $1;