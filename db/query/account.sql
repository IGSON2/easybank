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
WHERE owner = ? AND currency = ?;

-- name: GetAccountByID :one
SELECT * FROM accounts
WHERE id=?;

-- name: GetCurrency :many
SELECT currency FROM accounts
WHERE owner = ?;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateAccount :execresult
UPDATE accounts SET balance = ? WHERE id = ?;

-- name: AddAccountBalance :execresult
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id);

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;