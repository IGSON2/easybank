package db

import (
	"context"
	"database/sql"
	"fmt"
)

type runTx func(*Queries) error

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// [Composer]
// SqlStore provides all functions to execute db queries and transaction
type SqlStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}

// ExecTx executes a function within a database transaction
func (s *SqlStore) execTx(ctx context.Context, fn runTx) error {
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

var txKey = struct{}{}

// TransferTx perfroms a money transfer from one account to the other.
// It creates the transfer, add account entries, and update account's balance within a database transaction
func (s *SqlStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	// [Closer]
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var createErr error
		var id int64
		var creatResult sql.Result

		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")

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

		fmt.Println(txName, "create entry 1")
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

		fmt.Println(txName, "create entry 2")
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

		fmt.Println(txName, "get account 1")
		fromAcnt, getAcntErr := q.GetAccount(ctx, arg.FromAccountID)
		if getAcntErr != nil {
			return fmt.Errorf("can't get to from_account by id. err : %v", getAcntErr)
		}

		fmt.Println(txName, "get account 2")
		toAcnt, getAcntErr := q.GetAccount(ctx, arg.ToAccountID)
		if getAcntErr != nil {
			return fmt.Errorf("can't get to to_account by id. err : %v", getAcntErr)
		}

		if arg.FromAccountID < arg.ToAccountID {

			fmt.Println(txName, "update account1")
			var updateErr error
			_, updateErr = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: fromAcnt.Balance - arg.Amount,
			})
			if updateErr != nil {
				return fmt.Errorf("can't update account1. err : %v", updateErr)
			}

			fmt.Println(txName, "update account2")
			_, updateErr = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: toAcnt.Balance + arg.Amount,
			})
			if updateErr != nil {
				return fmt.Errorf("can't update account2. err : %v", updateErr)
			}
		} else {
			var updateErr error
			fmt.Println(txName, "update account2")
			_, updateErr = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: toAcnt.Balance + arg.Amount,
			})
			if updateErr != nil {
				return fmt.Errorf("can't update account2. err : %v", updateErr)
			}

			fmt.Println(txName, "update account1")
			_, updateErr = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: fromAcnt.Balance - arg.Amount,
			})
			if updateErr != nil {
				return fmt.Errorf("can't update account1. err : %v", updateErr)
			}
		}

		result.FromAccount, getAcntErr = q.GetAccount(ctx, arg.FromAccountID)
		if getAcntErr != nil {
			return fmt.Errorf("can't get to from_account by id. err : %v", getAcntErr)
		}

		result.ToAccount, getAcntErr = q.GetAccount(ctx, arg.ToAccountID)
		if getAcntErr != nil {
			return fmt.Errorf("can't get to to_account by id. err : %v", getAcntErr)
		}

		return createErr
	})
	return result, err
}
