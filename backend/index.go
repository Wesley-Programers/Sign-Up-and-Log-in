package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)


type Data struct {
	name string
	email string
	password int
}

var dataSlice []Data

func handler(w http.ResponseWriter, r * http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	if r.Method == http.MethodPost {
		w.Header().Set("Content-Type", "text/plain: charset=utf-8")

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "ERROR: ", http.StatusBadRequest)
			return
		}


		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		passwordString := strconv.Itoa(password)

		newUsers := Data{name: name, email: email, password: password}
		dataSlice = append(dataSlice, newUsers)

	} else {
		http.Error(w, "METHOD NOT PERMITED", http.StatusMethodNotAllowed)	
	}
			
}

func sqlTable(db *sql.DB, nameData, emailData, passwordData string) {

	if db == nil {
		log.Fatal("ERROR ON SQL TABLE") 
	}
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS usuarios (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		email VARCHAR(100),
		password VARCHAR(1000),
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)

	if err != nil {
		log.Fatal("ERROR TRYING TO CREATE THE TABLE ", err)
	}
	
}

func main() {
	http.HandleFunc("/", handler)	
	fmt.Println("SERVER OPEN WITH GOLANG")
}
