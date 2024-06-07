package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

type ConnDB struct {
	Db *sqlx.DB
}

var db *ConnDB

func NewConnection() *ConnDB {
	log.Println("-----")
	log.Println(os.Getenv("GOOSE_DBSTRING"))
	log.Println("-----")
	db = &ConnDB{
		Db: sqlx.MustOpen("sqlite3", os.Getenv("GOOSE_DBSTRING")),
	}
	return db
}

func GetConnection() *ConnDB {

	db = NewConnection()

	return db
}

func CloseConnection() {
	db.Db.Close()
}
