package database

import (
	"chadrss/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	// Connect to the database
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection is successful
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	DB = db
	fmt.Println("😁 Connected to database")
}
