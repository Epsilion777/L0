package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	username = "postgres"
	password = "123321"
	hostname = "localhost"
	port     = "5436"
	dbname   = "wb"
)

// postgres://postgres:123321@localhost:5436/wb?sslmode=disable
func InitPostgres(_ context.Context) *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, hostname, port, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
