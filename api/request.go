package api

import "fmt"

type createAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}

func (c *createAccountRequest) ValidateValues() error {
	if c.Currency == "" || c.Owner == "" {
		return fmt.Errorf("owner and currency must be required")
	}

	currencies := []string{"USD", "KRW", "EUR", "JAP", "BTC", "ETH"}
	var isSame bool
	for _, cur := range currencies {
		if c.Currency == cur {
			isSame = true
		}
	}
	if !isSame {
		return fmt.Errorf("invalid currency type %v", c.Currency)
	}
	return nil
}

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
