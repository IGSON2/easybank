package api

import (
	"database/sql"
	db "easybank/db/sqlc"
	"easybank/token"
	"easybank/util"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" validator:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validator:"required,min=1"`
	Amount        int64  `json:"amount" validator:"required,gt=0"`
	Currency      string `json:"currency" validator:"required,currency"`
}

func (server *Server) createTransfer(c *fiber.Ctx) error {
	var req transferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{err})
	}

	errs := util.ValidateStruct(req)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{errs})
	}

	fromaccount, isValid := server.validAccount(c, req.FromAccountID, req.Currency)
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{errors.New("invalid from account")})
	}

	authPayload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{errors.New("invalid payload")})
	}

	if fromaccount.Owner != authPayload.Username {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{errors.New("from account does not belong to the user")})
	}

	_, isValid = server.validAccount(c, req.ToAccountID, req.Currency)
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{errors.New("invalid to account")})
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(c.Context(), arg)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{err})
	}

	return c.Status(http.StatusOK).JSON(result)
}

func (s *Server) validAccount(c *fiber.Ctx, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccountByID(c.Context(), accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(fiber.StatusNotFound).JSON(errorResponse{err})
			return account, false
		}
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
		return account, false
	}

	if account.Currency != currency {
		err := errors.New("invalid currency")
		c.Status(fiber.StatusBadRequest).JSON(errorResponse{err})
		return account, false
	}

	return account, true

}
