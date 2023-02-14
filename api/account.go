package api

import (
	"database/sql"
	db "easybank/db/sqlc"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"
)

type errorResponse struct {
	RepErr error `json:"error"`
}

func (s *Server) createAccount(c *fiber.Ctx) error {
	var req createAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("cannot parse request body. err : %v", err)
	}
	if err := req.ValidateValues(); err != nil {
		c.Status(fiber.StatusBadRequest).SendString("Bad Request")
		return err
	}
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}
	result, err := s.store.CreateAccount(c.Context(), arg)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}
	return c.Status(fiber.StatusOK).JSON(lastID)
}

func (s *Server) getAccount(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{err})
	}
	if id <= 0 {
		return fmt.Errorf("id must be higher than 1")
	}
	account, err := s.store.GetAccount(c.Context(), int64(id))

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse{err})
		}
		fmt.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}
	return c.Status(fiber.StatusOK).JSON(account)
}

func (s *Server) listAccount(c *fiber.Ctx) error {
	var req listAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("cannot parse request body. err : %v", err)
	}
	if err := req.ValidateValues(); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse{err})
	}
	accounts, err := s.store.ListAccounts(c.Context(), db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}
	return c.Status(fiber.StatusOK).JSON(accounts)
}
