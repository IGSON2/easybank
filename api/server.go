package api

import (
	db "easybank/db/sqlc"

	"github.com/gofiber/fiber/v2"
)

// Server serves HTTP requests for our backing service.
type Server struct {
	store  *db.Store
	router fiber.App
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := fiber.New()
	router.Get("/", server.createAccount)
	return server
}
