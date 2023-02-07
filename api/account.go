package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type createAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency"`
}

func (s *Server) createAccount(c *fiber.Ctx) error {
	var req createAccountRequest
	// TODO : Validate value which is correct or not
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("Cannot parse request body. err : %v", err)
	}
	if err := req.validateValues(); err != nil {
		c.Status(fiber.StatusBadRequest).SendString("Bad Request")
		return err
	}

	// 나머지 완셩하기

}

func (c *createAccountRequest) validateValues() error {
	currencies := []string{"USD", "KRW", "EUR", "JAP", "BTC", "ETH"}
	for _, cur := range currencies {
		if c.Currency != cur {
			return fmt.Errorf("Invalid currency type %v", c.Currency)
		}
	}
}
