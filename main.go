package main

import (
	"database/sql"
	"easybank/api"
	db "easybank/db/sqlc"
	"easybank/util"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	util.LoadConfig(".")
}

func main() {
	conn, err := sql.Open(util.C.DBDriver, util.C.DBsource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(&store)
	log.Fatalln(server.Start(util.C.Port))
}
