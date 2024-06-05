package db

import (
	"os"

	"github.com/jmoiron/sqlx"
)

type  ConnDB struct{
	Db *sqlx.DB
}

var db *ConnDB

func NewConnection() *ConnDB {
	db = &ConnDB{
		Db: sqlx.MustOpen("sqlite3", os.Getenv("GOOSE_DBSTRING")),
	}
	return db
}

func GetConnection() *ConnDB {
	if db == nil{
		db=NewConnection()
	}
	return db
}

