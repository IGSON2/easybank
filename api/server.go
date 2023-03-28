package api

import (
	db "easybank/db/sqlc"
	"easybank/token"
	"easybank/util"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Server serves HTTP requests for our backing service.
type Server struct {
	config     util.Config
	store      db.Store
	router     *fiber.App
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}
	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: maker,
	}
	router := fiber.New()
	router.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	router.Post("/account", server.createAccount)
	router.Get("/account", server.listAccount)
	router.Get("/account/:id", server.getAccount)
	router.Post("/transfer", server.createTransfer)
	router.Post("/user", server.createUser)
	router.Post("/user/login", server.loginUser)
	server.router = router
	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Listen(address)
}
