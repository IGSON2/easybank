package api

import (
	db "easybank/db/sqlc"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (server *Server) goodAccountCurrency(c *fiber.Ctx, accountID int64, currency string) error {
	account, err := server.store.GetAccountByID(c.Context(), accountID)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(errorResponse{err})
		return err
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		c.Status(http.StatusBadRequest).JSON(errorResponse{err})
		return err
	}

	return nil
}

func (server *Server) createTransfer(c *fiber.Ctx) error {
	var req transferRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(http.StatusBadRequest).JSON(errorResponse{err})
		return err
	}

	if err := server.goodAccountCurrency(c, req.FromAccountID, req.Currency); err != nil {
		return err
	}

	if err := server.goodAccountCurrency(c, req.ToAccountID, req.Currency); err != nil {
		return err
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(c.Context(), arg)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(errorResponse{err})
		return err
	}

	c.Status(http.StatusInternalServerError).JSON(result)
	return nil
}
