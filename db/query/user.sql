-- name: CreateUser :execresult
INSERT INTO
    users (
        username,
        hashed_password,
        full_name,
        email
    )
VALUES
    (?, ?, ?, ?);

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    username = ?
LIMIT
    1;

-- name: GetLastUser :one
SELECT
    *
FROM
    users
ORDER BY
    created_at DESC
LIMIT
    1;