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
	CreateTableResetPassword(database)
	CreateTableAttempts(database)
	return database
}


func CreateTable(database *sql.DB) {

	if database == nil {
		log.Fatal("ERROR") 
	}

	_, err := database.Exec(`
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


func CreateTableResetPassword(database *sql.DB) {

	if database == nil {
		log.Fatal("ERROR")
	}

	_, err := database.Exec(`
	CREATE TABLE IF NOT EXISTS reset_password (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		token VARCHAR(255) UNIQUE NOT NULL,
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCATE
	)`)

	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}


func CreateTableAttempts(database *sql.DB) {

	if database == nil {
		log.Fatal("ERROR")
	}

	_, err := database.Exec(`
	CREATE TABLE IF NOT EXISTS attempts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(120) NOT NULL,
		attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX index_email_time(email, attempted_at)
	)`)

	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}
