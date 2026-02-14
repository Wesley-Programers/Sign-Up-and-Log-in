package database

import (
	"fmt"
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


func Connect() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	// log.SetFlags(log.Lshortfile)

	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")

	dataSourceTime := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", databaseUser, databasePassword, databaseHost, databasePort, databaseName)

	database, err := sql.Open("mysql", dataSourceTime) 
	if err != nil {
		log.Fatal("ERRO: ", err)
	}

	if err = database.Ping(); err != nil {
        log.Fatal("ERROR: ", err)
    }

	CreateTable(database)
	return database

}


func CreateTable(db *sql.DB) {

	// log.SetFlags(log.Lshortfile)

	if db == nil {
		log.Fatal("ERROR ON SQL TABLE") 
	}

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(120) UNIQUE NOT NULL,
		password VARCHAR(150) NOT NULL,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		log.Fatal("ERROR: ", err)
	}

}
