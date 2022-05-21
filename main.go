package main

import (
	"database/sql"
	"log"

	db "github.com/techschool/simplebank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/sds?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Connot connect to DB: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start()
}
