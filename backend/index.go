package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)


type Data struct {
	name string
	email string
	password string
}

type Claims struct {
	Name string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("")

var dataSlice []Data


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

			nameOrEmail := r.FormValue("nameEmail")
			passwordLog := r.FormValue("passwordLog")

			err2 := verifyLogin(database, nameOrEmail, nameOrEmail, passwordLog)

			if err2 == nil {
				
				fmt.Println("")
				log.Println("")
				w.WriteHeader(200)
				w.Write([]byte(""))

			} else if err2.Error() == "nome nao existe" {
				log.Println("")
				w.WriteHeader(409)
				w.Write([]byte(""))
				fmt.Println("")

			} else if err2.Error() == "senha nao existe" {
				log.Println("")
				w.WriteHeader(409)
				w.Write([]byte(""))
				fmt.Println("")

			} else if err2.Error() == "email nao existe" {
				log.Println("")
				w.WriteHeader(409)
				w.Write([]byte(""))
				fmt.Println("")

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


func sqlTable(db *sql.DB, nameData, emailData, passwordData string) {

	if db == nil {
		log.Fatal("") 
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
		log.Fatal("", err)
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


func verifyLogin(database *sql.DB, nameLogin, passwordLogin string) bool {
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
		return errors.New("")

	}
	fmt.Println("")
	return nil
}


func generateToken(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/plain")


	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")

	expiration := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Name: name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("ERROR: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: stringToken,
		Expires: expiration,
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteStrictMode,
	})

	log.Println("")
}


func validToken(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token invalid or expired", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"name": claims.Name,
		"email": claims.Email,
	}
	json.NewEncoder(w).Encode(response)
}


func main() {
	database := database()
	defer database.Close()

	http.HandleFunc("/sign", handler(database))
	http.HandleFunc("/login", handlerLogIn(database))
	
	fmt.Println("SERVER OPEN WITH GOLANG")
}
