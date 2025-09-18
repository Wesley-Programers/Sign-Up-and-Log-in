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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

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

		if newUsers.name != "" && newUsers.email != "" && newUsers.password != "" {
			
			dataSlice = append(dataSlice, newUsers)
			fmt.Printf("\nName: %v\nemail: %v\npassword: %v\n", newUsers.name, newUsers.email, newUsers.password)
			fmt.Println(dataSlice)
			// fmt.Printf("%v\n", len(dataSlice))
			nameData := newUsers.name
			emailData := newUsers.email
			passwordData := newUsers.password
			database(nameData, emailData, passwordData)

			fmt.Println("NAME DUPLICATE ON HANDLER FUNC: ", nameDuplicate)
			fmt.Println("VERIFY HELP ON HANDLER FUNC: ", verifyHelp)


			if verifyHelp {
				if nameDuplicate {
					w.WriteHeader(200)
					fmt.Println("ERROR 202")
					return
				}
			} else if !verifyHelp {
				http.Redirect(w, r, "", http.StatusFound)
				return
			}

		} else {
			fmt.Println("\nVALORES VAZIOS NO HANDLER")
		}

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

func sqlInsert(databasePointer *sql.DB, nameData, emailData, passwordData string) error {

	query := "INSERT INTO usuarios (name, email, password) VALUES (?, ?, ?);"
		
	if nameData != ""  && emailData != "" && passwordData != "" {
		_, erroInsert := databasePointer.Exec(query, nameData, emailData, passwordData)

		if erroInsert != nil {
			fmt.Println("ERROR TRYING TO INSERT THE DATA ", erroInsert)

			if strings.Contains(erroInsert.Error(), "Error 1062") {

				verifyHelp = true
				if strings.Contains(erroInsert.Error(), "for key 'unique_name'") {

					fmt.Printf("LAST NAME ALREADY EXISTS: %s\n", nameData)
					nameDuplicate = true

				} else if strings.Contains(erroInsert.Error(), "for key 'unique_email'") {

					fmt.Printf("LAST EMAIL ALREADY EXISTS: %s\n", emailData)
					emailDuplicate = true

				} // IF STRING CONTAINS "FOR KEY 'UNIQUE_NAME"
			} else {

				verifyHelp = false
				nameDuplicate = false
				emailDuplicate = false

			}// IF STRINGS.CONTAINS "ERROR 1062"
		} // IF ERRO INSERT != NIL
	} else {
		fmt.Println("OS VALORES ESTÃO VAZIOS , PRINT DA FUNÇÃO SQL INSERT\n")
	}

	return nil
}

func main() {
	http.HandleFunc("/", handler)	
	fmt.Println("SERVER OPEN WITH GOLANG")
}
