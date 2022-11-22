package db

import (
	"context"
	"database/sql"
	"fmt"
)

type runTx func(*Queries) error

// [Composer]
// Store provides all functions to execute db queries and transaction
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return Store{
		Queries: New(db),
		db:      db,
	}
}

// ExecTx executes a function within a database transaction
func (s *Store) execTx(ctx context.Context, fn runTx) error {
	// BeginTx 의 option 인자로 격리 레벨을 조정할 수 있다.
	Tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(Tx)

	err = fn(q)
	if err != nil {
		if rbErr := Tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, RollbackErr : %v", err, rbErr)
		}
		return err
	}

	// Tx가 에러 없이 성공했을 경우
	return Tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx perfroms a money transfer from one account to the other.
// It creates the transfer, add account entries, and update account's balance within a database transaction
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	// [Closer]
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var createErr error
		var id int64
		var creatResult sql.Result

		// Transfer
		creatResult, createErr = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if createErr != nil {
			return fmt.Errorf("can't create transfer. err : %v", createErr)
		}
		id, createErr = creatResult.LastInsertId()
		if createErr != nil {
			return fmt.Errorf("can't get id by created transfer. err : %v", createErr)
		}
		result.Transfer, createErr = q.GetTransfer(ctx, id)
		if createErr != nil {
			return fmt.Errorf("can't get transfer by id. err : %v", createErr)
		}

		// FromEntry
		creatResult, createErr = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if createErr != nil {
			return fmt.Errorf("can't create from_entry. err : %v", createErr)
		}
		id, createErr = creatResult.LastInsertId()
		if createErr != nil {
			return fmt.Errorf("can't get id by created from entry. err : %v", createErr)
		}
		result.FromEntry, createErr = q.GetEntry(ctx, id)
		if createErr != nil {
			return fmt.Errorf("can't get from entry by id. err : %v", createErr)
		}

		// ToEntry
		creatResult, createErr = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if createErr != nil {
			return fmt.Errorf("can't create to_entry. err : %v", createErr)
		}
		id, createErr = creatResult.LastInsertId()
		if createErr != nil {
			return fmt.Errorf("can't get id by created to entry. err : %v", createErr)
		}
		result.ToEntry, createErr = q.GetEntry(ctx, id)
		if createErr != nil {
			return fmt.Errorf("can't get to entry by id. err : %v", createErr)
		}

		// FromAccount
		_, createErr = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if createErr != nil {
			return fmt.Errorf("can't create from_account. err : %v", createErr)
		}

		result.FromAccount, createErr = q.GetAccount(ctx, arg.FromAccountID)
		if createErr != nil {
			return fmt.Errorf("can't get to from_account by id. err : %v", createErr)
		}

		// ToAccount
		creatResult, createErr = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if createErr != nil {
			return fmt.Errorf("can't create to_account. err : %v", createErr)
		}

		result.ToAccount, createErr = q.GetAccount(ctx, arg.ToAccountID)
		if createErr != nil {
			return fmt.Errorf("can't get to to_account by id. err : %v", createErr)
		}

		return createErr
	})
	return result, err
}
