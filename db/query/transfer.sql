-- name: CreateTransfer :execresult
INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
) VALUES (
    ?, ?, ?
);

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ? LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE
    from_account_id = ? OR
    to_account_id = ?
ORDER BY id
LIMIT ?
OFFSET ?;