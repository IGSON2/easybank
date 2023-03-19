package db

import (
	"database/sql"
	"easybank/util"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	util.LoadConfig("../../.")
	var err error

	testDB, err = sql.Open(util.C.DBDriver, util.C.DBsource)
	if err != nil {
		log.Fatal("Can't open th db : ", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
