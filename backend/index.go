package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)


type Data struct {
	name string
	email string
	password string
}

var dataSlice []Data

var verifyHelp bool = false
var nameDuplicate bool = false
var emailDuplicate bool = false

func handler(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/plain")


		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {		

			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				http.Error(w, "ERROR: ", http.StatusBadRequest)
				return
			}


			name := r.FormValue("name")
			email := r.FormValue("email")
			password := r.FormValue("password")
			hash, err := hashPassword(password)
			
			if err != nil {
				log.Fatal("ERROR ON HASH: ", err)
			}
			// passwordString := strconv.Itoa(password)
			
			newUsers := Data{name: name, email: email, password: hash}
			
			if newUsers.name != "" && newUsers.email != "" && newUsers.password != "" {
				
				fmt.Printf("\nName: %v\nemail: %v\npassword: %v\n", newUsers.name, newUsers.email, newUsers.password)
				fmt.Println(dataSlice)

				nameData := newUsers.name
				emailData := newUsers.email
				passwordData := newUsers.password

				errInsert := sqlInsert(database, nameData, emailData, passwordData)
				if errInsert != nil {
					http.Error(w, "Some error", http.StatusInternalServerError)
					return
				}
				
				dataSlice = append(dataSlice, newUsers)

				fmt.Println("NAME DUPLICATE ON HANDLER FUNC: ", nameDuplicate)
				fmt.Println("EMAIL DUPLICATE ON HANDLER FUNC: ", emailDuplicate)
				fmt.Println("VERIFY HELP ON HANDLER FUNC: ", verifyHelp)


				if verifyHelp {
					verifyHelp = false
					
					if nameDuplicate {
						log.Println("Mandando o codigo 409 No nameDuplicate")
						w.WriteHeader(409)
						w.Write([]byte("Nome já existe"))
						nameDuplicate = false
						return
						
					} else if emailDuplicate {
						log.Println("Mandando o codigo 409 No emailDuplicate")
						w.WriteHeader(409)
						w.Write([]byte("Email já existe"))
						emailDuplicate = false
						return
					}
					
				} else {
					
					log.Println("Mandando o codigo 201")
					w.WriteHeader(201)
					w.Write([]byte("Dados validos"))
					return

				}

			} else {
				// w.WriteHeader(http.StatusOK)
				fmt.Println("\nVALORES VAZIOS NO HANDLER")
			}
		
			
		} else {
			http.Error(w, "METHOD NOT PERMITED", http.StatusMethodNotAllowed)	
		}

	}

}


func handlerLogIn(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/plain")


		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {

			err := r.ParseForm()
			if err != nil {
				http.Error(w, "ERROR: ", http.StatusBadRequest)
				return
			}

			nameOrEmail := r.FormValue("nameEmail")
			passwordLog := r.FormValue("passwordLog")
			fmt.Println(nameOrEmail)
			if verifyLogin(database, nameOrEmail, passwordLog) {
				fmt.Println("SIM, LOGIN FEITO")
			} else {
				fmt.Println("NAO, LOGIN NAO FEITO")
			}
			
		} else {
			http.Error(w, "METHOD NOT PERMITED", http.StatusMethodNotAllowed)
		}

	}
}

func database() *sql.DB {
	
}


func sqlTable(db *sql.DB, nameData, emailData, passwordData string) {

	if db == nil {
		log.Fatal("ERROR ON SQL TABLE WHERE *SQL.DB == NIL") 
	}

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS usuarios (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(80),
		email VARCHAR(80),
		password VARCHAR(100),
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

				}
			} else {

				verifyHelp = false
				nameDuplicate = false
				emailDuplicate = false

			}
		}
	} else {
		fmt.Println("OS VALORES ESTÃO VAZIOS , PRINT DA FUNÇÃO SQL INSERT")
	}

	return nil
}


func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


func verifyLogin(database *sql.DB, nameLogin, passwordLogin string) bool {
	query := "SELECT password FROM usuarios WHERE name = ?";

	var passwordHash string
	err := database.QueryRow(query, nameLogin).Scan(&passwordHash)
	if err == sql.ErrNoRows {
		fmt.Println("\nNOME NAO EXISTE")
		return false
	} else if err != nil {
		log.Fatal("ERROR NA FUNÇÃO VERIFY LOGIN: ", err)
		return false
	}

	errPasswordHash := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordLogin))

	if errPasswordHash != nil {
		fmt.Println("SENHA INVALIDA")
		return false
	}
	
	fmt.Println("NOME E TALVEZ A SENHA EXISTE SIM")
	return true
}


func main() {
	database := database()
	defer database.Close()

	http.HandleFunc("/sign", handler(database))
	http.HandleFunc("/login", handlerLogIn(database))
	
	fmt.Println("SERVER OPEN WITH GOLANG")
}
