package api

import (
	db "easybank/db/sqlc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Server serves HTTP requests for our backing service.
type Server struct {
	store  *db.Store
	router *fiber.App
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := fiber.New()
	router.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	router.Post("/", server.createAccount)
	router.Get("/accounts/:id", server.getAccount)
	router.Get("/accounts", server.listAccount)
	server.router = router
	return server
}

func (s *Server) Start(address string) error {
	return s.router.Listen(address)
}
