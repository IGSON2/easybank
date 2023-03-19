package api

import (
	db "easybank/db/sqlc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Server serves HTTP requests for our backing service.
type Server struct {
	store  db.Store
	router *fiber.App
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := fiber.New()
	router.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	router.Post("/account", server.createAccount)
	router.Get("/account/:id", server.getAccount)
	router.Get("/account", server.listAccount)
	router.Post("/transfer", server.createTransfer)
	router.Post("/user", server.createUser)
	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Listen(address)
}
