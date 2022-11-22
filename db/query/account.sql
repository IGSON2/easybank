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

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1 FOR NO KEY UPDATE; -- TODO : Avoid Deadlock in mySQL

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