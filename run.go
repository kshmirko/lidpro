package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kshmirko/lidpro/db"
	_ "github.com/mattn/go-sqlite3"
)


func main(){
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal("Ошибка загрузки файла .env!")
	}
	
	con:=db.GetConnection()
	con.Db.Ping()
	con.Db.Close()
}