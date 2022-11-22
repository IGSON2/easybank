package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriverName = "mysql"
	dbSourceName = "root:123@tcp(localhost:3306)/easy_bank?parseTime=true"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriverName, dbSourceName)
	if err != nil {
		log.Fatal("Can't open th db : ", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
