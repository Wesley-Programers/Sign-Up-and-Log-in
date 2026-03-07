package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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
	databaseConfig := os.Getenv("DATABASE_CONFIG")

	dataSourceTime := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", databaseUser, databasePassword, databaseHost, databasePort, databaseName, databaseConfig)

	database, err := sql.Open("mysql", dataSourceTime) 
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	if err = database.Ping(); err != nil {
        log.Fatal("ERROR: ", err)
    }

	database.SetMaxOpenConns(20)
	database.SetMaxIdleConns(10)
	database.SetConnMaxLifetime(10 * time.Minute)

	createTable(database)
	createTableResetPassword(database)
	createTableAttempts(database)
	createLoginAttempts(database)
	
	return database
}


func createTable(database *sql.DB) {

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


func createTableResetPassword(database *sql.DB) {

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
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`)

	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}


func createTableAttempts(database *sql.DB) {

	if database == nil {
		log.Fatal("ERROR")
	}

	_, err := database.Exec(`
	CREATE TABLE IF NOT EXISTS attempts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(120) NOT NULL,
		attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX index_email_time(email, attempted_at)
	);`)

	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}

func createLoginAttempts(database *sql.DB) {
	
	if database == nil {
		log.Fatal("ERROR")
	}

	_, err := database.Exec(`
	CREATE TABLE IF NOT EXISTS login_attempts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(120) NOT NULL,
		attempt_in TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		success TINYINT(1) DEFAULT 0,
		INDEX index_email(email, attempt_in)
	);`)

	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}
