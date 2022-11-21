-- name: CreateAccount :execresult
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  ?,?,?
);

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateAccount :execresult
UPDATE accounts SET balance = ? WHERE id = ?;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;