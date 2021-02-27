package main

import (
	"log"

	"github.com/FelipeAz/desafio-serasa/app/infrastructure"
)

func main() {
	db, err := infrastructure.NewSQLHandler()
	if err != nil {
		log.Fatal(err)
	}

	defer db.CloseConnection()
}