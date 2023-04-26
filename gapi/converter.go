package gapi

import (
	db "easybank/db/sqlc"
	"easybank/pb"
)

func convertUser(dbUser *db.User) *pb.User {
	return &pb.User{
		Username: dbUser.Username,
		FullName: dbUser.FullName,
		Email:    dbUser.Email,
	}
}
