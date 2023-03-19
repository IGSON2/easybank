package api

import (
	db "easybank/db/sqlc"
	"easybank/util"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,alphanum,max=12"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createUser(ctx *fiber.Ctx) error {
	var req createUserRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(err)
		return err
	}

	errors := util.ValidateStruct(req)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse{err})
		return err
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	_, err = s.store.CreateUser(ctx.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.Status(http.StatusForbidden).JSON(errorResponse{err})
				return err
			}
		}
		ctx.Status(http.StatusInternalServerError).JSON(errorResponse{err})
		return err
	}

	user, err := s.store.GetUser(ctx.Context(), arg.Username)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(err)
		return err
	}
	rsp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.Status(http.StatusOK).JSON(rsp)
	return nil
}
