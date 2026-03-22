package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	var database *sql.DB

	for i := range 10 {
		database, err = sql.Open("mysql", dataSourceTime)
		if err == nil {
			err = database.Ping()
			if err == nil {
				log.Println("SUCCESS")

				database.SetMaxOpenConns(20)
				database.SetMaxIdleConns(10)
				database.SetConnMaxLifetime(10 * time.Minute)
				return database
			}
		}

		log.Printf("WAITING... ATTEMPT: %d/10", i+10)
		time.Sleep(3 * time.Second)
	}

	log.Fatal("ERROR: ", err)
	return nil
}


func RunMigrations(db *sql.DB) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("ERROR: ", err)
	}

	log.Println("SUCCESS")
}
