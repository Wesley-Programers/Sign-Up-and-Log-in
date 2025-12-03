package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	
	// "encoding/json"
	// "strconv"
	// "io"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	// "github.com/golang-jwt/jwt/v5"
)


type Data struct {
	name string
	email string
	password string
}

type Session struct {
	UserID int
	Expires time.Time
}

var dataSlice []Data


func RandString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}


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
			
			newUsers := Data{name: name, email: email, password: hash}
			
			if newUsers.name != "" && newUsers.email != "" && newUsers.password != "" {
				
				fmt.Printf("\nName: %v\nemail: %v\npassword: %v\n", newUsers.name, newUsers.email, newUsers.password)
				fmt.Println(dataSlice)

				nameData := newUsers.name
				emailData := newUsers.email
				passwordData := newUsers.password

				errInsert := sqlInsert(database, nameData, emailData, passwordData)

				if errInsert.Error() == "dados validos" {
					dataSlice = append(dataSlice, newUsers)
					log.Println("")
					w.WriteHeader(201)
					w.Write([]byte(""))
					dataSlice = append(dataSlice, newUsers)

				} else if errInsert.Error() == "nome ja existe" {
					log.Println("")
					w.WriteHeader(409)
					w.Write([]byte(""))

				} else if errInsert.Error() == "email ja existe" {
					log.Println("")
					w.WriteHeader(409)
					w.Write([]byte(""))

				} else if errInsert != nil {
					http.Error(w, "Some error", http.StatusInternalServerError)

				} else {
					log.Println("")
				}

			} else {
				// w.WriteHeader(http.StatusOK)
				fmt.Println("")
			}
		
			
		} else {
			http.Error(w, "METHOD NOT PERMITED", http.StatusMethodNotAllowed)	
		}

	}

}


func handlerLogIn(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
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

			nameOrEmail := r.FormValue("nameEmail")
			passwordLog := r.FormValue("passwordLog")

			err2 := verifyLogin(database, nameOrEmail, nameOrEmail, passwordLog)

			if err2 == nil {
				fmt.Println("\nSIM, LOGIN FEITO")
				log.Println("MANDANDO O CODIGO 200 NO VERIFY LOGIN")
				w.WriteHeader(200)
				w.Write([]byte("Dados de login validos"))

				sessionID := RandString(64)
				var userID int

				var passwordHashed string

				anotherErr := database.QueryRow(
				"SELECT id, password FROM usuarios WHERE email = ?", nameOrEmail).Scan(&userID, &passwordHashed)
				if anotherErr != nil {
					log.Fatalf("ERROR TRYIN' TO GET THE USER ID: %v", anotherErr)
				}

				errHash := bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(passwordLog))
				if errHash != nil {
					fmt.Println("senha errada")
					return
				}

				err := SaveSession(database, sessionID, userID, time.Now().Add(5 * time.Minute))
				if err != nil {
					http.Error(w, "ERROR", 500)
					return
				}

				cookie := http.Cookie{
					Name: "id_user",
					Value: sessionID,
					Path: "/",
					HttpOnly: false,
					Secure: false,
					SameSite: http.SameSiteStrictMode,
					Expires: time.Now().Add(5 * time.Minute),
				}
				
				http.SetCookie(w, &cookie)

			} else if err2.Error() == "nome nao existe" {
				log.Println("Mandando o codigo 409 no verify login nome nao existe")
				w.WriteHeader(409)
				w.Write([]byte("Nome incorreto"))
				fmt.Println("Nome nao existe")

			} else if err2.Error() == "senha nao existe" {
				log.Println("Mandando o codigo 409 no verify login senha nao existe")
				w.WriteHeader(409)
				w.Write([]byte("Senha incorreta"))
				fmt.Println("Senha nao existe")

			} else if err2.Error() == "email nao existe" {
				log.Println("Mandando o codigo 409 no verify login email nao existe")
				w.WriteHeader(409)
				w.Write([]byte("Email incorreto"))
				fmt.Println("Email nao existe")

			} else {
				log.Println("Error inesperado")
			}
			
		} else {
			http.Error(w, "METHOD NOT PERMITED", http.StatusMethodNotAllowed)
		}
	}
}


func database() *sql.DB {

	return database
}


func sqlTable(db *sql.DB) {

	if db == nil {
		log.Fatal("ERROR ON SQL TABLE") 
	}

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS usuarios (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(80),
		email VARCHAR(80) UNIQUE,
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

				if strings.Contains(erroInsert.Error(), "for key 'unique_name'") {

					fmt.Printf("LAST NAME ALREADY EXISTS: %s\n", nameData)
					return errors.New("")

				} else if strings.Contains(erroInsert.Error(), "for key 'unique_email'") {

					fmt.Printf("LAST EMAIL ALREADY EXISTS: %s\n", emailData)
					return errors.New("")

				}
			} else {
				return errors.New("")
			}
		}
	} else {
		fmt.Println("")
	}

	return nil
}


func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


func verifyLogin(database *sql.DB, nameLogin, emailLogin, passwordLogin string) (error) {
	query := "SELECT password FROM usuarios WHERE name = ? OR email = ?";

	var passwordHash string
	err := database.QueryRow(query, nameLogin, emailLogin).Scan(&passwordHash)
	
	if err == sql.ErrNoRows {
		if strings.HasSuffix(nameLogin, "") || strings.HasSuffix(emailLogin, "") {
			return errors.New("")
		} else {
			return errors.New("")
		}
	} else if err != nil {
		log.Fatal("some error: ", err)
	}

	
	errPasswordHash := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordLogin))

	if errPasswordHash != nil {
		fmt.Println("")
		return errors.New("")
	}
	
	fmt.Println("")
	return nil
}


func main() {
	database := database()
	defer database.Close()

	http.HandleFunc("/sign", handler(database))
	http.HandleFunc("/login", handlerLogIn(database))
	http.HandleFunc("/main", handlerRecoverCookie(database))
	
	fmt.Println("SERVER OPEN WITH GOLANG")
}
