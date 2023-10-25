package queries

import (
	"database/sql"
	"log"
	"os"
)

var Conn *sql.DB

func Connect() {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	Conn = db
}

func Disconnect() {
	Conn.Close()
}

var Q *Queries

func Build() {
	Q = New(Conn)
}
