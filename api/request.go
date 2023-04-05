package api

import (
	"easybank/util"
	"fmt"
)

type listAccountRequest struct {
	PageID   int32 `json:"pageid"`
	PageSize int32 `json:"pagesize"`
}

func (l *listAccountRequest) ValidateValues() error {
	if l.PageID <= 0 {
		return fmt.Errorf("page_id must be higher than 0")
	}

	if l.PageSize < 5 || l.PageSize > 10 {
		return fmt.Errorf("page_size must set between 5 to 10")
	}
	return nil
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

func (t *transferRequest) ValidateValues() error {
	if t.FromAccountID <= 0 || t.ToAccountID <= 0 {
		return fmt.Errorf("AccountID must be higher than 0 From : %d, To : %d", t.FromAccountID, t.ToAccountID)
	}

	if t.Amount <= 0 {
		return fmt.Errorf("amount must be higher than 0")
	}

	var contain bool
	for _, s := range util.Currencies {
		if t.Currency == s {
			contain = true
		}
	}
	if !contain {
		return fmt.Errorf("invalid currency type")
	}

	return nil
}
