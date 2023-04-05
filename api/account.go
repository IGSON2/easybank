package api

import (
	"database/sql"
	db "easybank/db/sqlc"
	"easybank/token"
	"easybank/util"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"
)

type errorResponse struct {
	RepErr error `json:"error"`
}

type createAccountRequest struct {
	Currency string `json:"currency" validate:"required,currency"`
}

func (s *Server) createAccount(c *fiber.Ctx) error {
	var req = createAccountRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("cannot parse request body. err : %v", err)
	}
	if errs := util.ValidateStruct(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{errs})
	}

	authPayload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{errors.New("invalid payload")})
	}

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	_, err := s.store.CreateAccount(c.Context(), arg)
	if err != nil {
		if msErr, ok := err.(*mysql.MySQLError); ok {
			switch msErr.Message {
			case "foreign_key_violation", "unique_violation":
				return c.Status(fiber.StatusForbidden).JSON(errorResponse{msErr})
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}
	return c.Status(fiber.StatusOK).JSON(req)
}

type getAccountRequest struct {
	ID int64 `json:"id" validate:"required,min=1"`
}

func (s *Server) getAccount(c *fiber.Ctx) error {
	var req getAccountRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{err})
	}

	errs := util.ValidateStruct(req)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{errs})
	}

	account, err := s.store.GetAccountByID(c.Context(), int64(req.ID))

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse{err})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}

	authPayload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{errors.New("invalid payload")})
	}

	if account.Owner != authPayload.Username {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{errors.New("account dosen't belong to the authenticated user")})
	}

	return c.Status(fiber.StatusOK).JSON(account)
}

type listAccountRequest struct {
	PageID   int32 `json:"pageid" validator:"required,min=1"`
	PageSize int32 `json:"pagesize" validator:"required,min=5,max=10"`
}

func (s *Server) listAccount(c *fiber.Ctx) error {
	var req listAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("cannot parse request body. err : %v", err)
	}

	errs := util.ValidateStruct(req)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{errs})
	}

	authPayload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{errors.New("invalid payload")})
	}

	accounts, err := s.store.ListAccounts(c.Context(), db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}

	return c.Status(fiber.StatusOK).JSON(accounts)
}
