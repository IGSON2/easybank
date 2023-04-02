package api

import (
	"database/sql"
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

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *fiber.Ctx) error {
	var req createUserRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(err)
		return err
	}

	errs := util.ValidateStruct(req)
	if errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
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
	rsp := newUserResponse(user)
	ctx.Status(http.StatusOK).JSON(rsp)
	return nil
}

type loginUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (s *Server) loginUser(ctx *fiber.Ctx) error {
	var req loginUserRequest

	err := ctx.BodyParser(&req)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(err)
		return err
	}

	errs := util.ValidateStruct(req)
	if errs != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errs)
	}

	user, err := s.store.GetUser(ctx.Context(), req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(errorResponse{err})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errorResponse{err})
	}

	// [TODO]
	// To use payload
	accessToken, _, err := s.tokenMaker.CreateToken(req.Username, s.config.AccessTokenDuration)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(errorResponse{err})
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	return ctx.Status(http.StatusOK).JSON(rsp)
}
