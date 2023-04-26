package gapi

import (
	db "easybank/db/sqlc"
	"easybank/pb"
	"easybank/token"
	"easybank/util"
	"fmt"
)

// Server serves HTTP requests for our backing service.
type Server struct {
	pb.UnimplementedEasybankServer
	config     util.Config
	store      db.Store
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
	return server, nil
}
