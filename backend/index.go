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

				if errInsert == nil {

					session, _ := store.Get(r, "tokenSession")
					session.Values["log"] = true
					session.Values["userName"] = newUsers.name
					session.Values["userEmail"] = newUsers.email
					session.Save(r, w)

					dataTest := fmt.Sprintf("%s|%s", newUsers.name,  newUsers.email)
					dataBase64 := base64.StdEncoding.EncodeToString([]byte(dataTest))

					cookie := &http.Cookie{
						Name: "user_data",
						Value: dataBase64,
						Path: "/",
						Expires: time.Now().Add(1 * time.Minute),
						HttpOnly: false,
						Secure: false,
						// SameSite: http.SameSiteLaxMode,
					}

					http.SetCookie(w, cookie)
					
					dataSlice = append(dataSlice, newUsers)
					log.Println("Mandando o codigo 201")
					w.WriteHeader(201)
					w.Write([]byte("Dados validos"))
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
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(80) NOT NULL,
		email VARCHAR(80) UNIQUE NOT NULL,
		password VARCHAR(159) NOT NULL,
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


func handlerDeleteAccount(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/plain")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.SetFlags(log.Lshortfile)

		if r.Method == http.MethodPost {

			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				http.Error(w, "ERROR: ", http.StatusBadRequest)
				return
			}

			emailConfirm := r.FormValue("emailConfirm")
			passwordConfirm := r.FormValue("passwordConfirm")

			var email bool
			var passwordHash string
			
			query := "DELETE FROM usuarios WHERE password = ?"
			queryEmailSelect := "SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = ?)"
			queryPasswordSelect := "SELECT password FROM usuarios WHERE email = ?"

			queryEmailErr := database.QueryRow(queryEmailSelect, emailConfirm).Scan(&email)
			if queryEmailErr != nil {
				log.Println("ERROR: ", queryEmailErr)
			}

			queryPasswordErr := database.QueryRow(queryPasswordSelect, emailConfirm).Scan(&passwordHash)
			if queryPasswordErr != nil {
				log.Println("ERROR: ", queryPasswordErr)

			} else if queryPasswordErr == sql.ErrNoRows {
				fmt.Println("SOMETHING")
				return
			}

			passwordHashCompare := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordConfirm))

			if passwordHashCompare == nil && email {
				fmt.Println("SENHA CORRETA")

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("EVERYTHING VALID"))
				_, err := database.Exec(query, passwordHashCompare)

				if err != nil {
					log.Fatal("SOMETHING BAD: ", err)
				}
				
			} else if passwordHashCompare != nil {
				fmt.Println("SENHA ERRADA")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("INCORRECT PASSWORD"))

			} else if !email {
				fmt.Println("INCORRECT EMAIL")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("INCORRECT EMAIl"))

			} else {
				log.Println("SOME ERROR")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("SOME ERROR"))
			}

		} else {
			http.Error(w, "ERROR: ", http.StatusMethodNotAllowed)
		}
	}
}


func resetPassword(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/plain")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.SetFlags(log.Lshortfile)

		if r.Method == http.MethodPost {

			var id int
			var verify string

			email := r.FormValue("email")

			query := "SELECT id, email FROM usuarios WHERE email = ?"

			queryError := database.QueryRow(query, email).Scan(&id, &verify)
			if queryError != nil {
				if queryError == sql.ErrNoRows {
					log.Println("EMAIL NOT FOUND")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("INVALID EMAIL"))	
					return
				}
				log.Println("ERROR: ", queryError)
				return
			}

			if verify == email {
				fmt.Println("VALID EMAIL, NEXT")

				teste = append(teste, email)
				fmt.Println("EMAIL IN SLICE HERE: ", teste[0])

				go func() {
					<-time.After(200 * time.Second)
					teste = teste[:0]
					fmt.Println("EMPTY SLICE")
				}()

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("VALID EMAIL"))

				token, err := generateTokens()
				if err != nil {
					log.Println("ERROR: ", err)
					return
				}

				expiresAt := time.Now().Add(5 * time.Minute)
				_, erro := database.Exec(
					"INSERT INTO password_token (user_id, token, expires_at) VALUES (?, ?, ?)", id, token, expiresAt,
				)
				if erro != nil {
					log.Println("ERROR: ", erro)
				}

				link := generateLink(token)
				w.Write([]byte(link))
				
			} else {
				log.Println("SOME ERROR")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("SOME ERROR"))
			}

		} else {
			http.Error(w, "METHOD NOT PERMITED", http.StatusMethodNotAllowed)
		}
	}
}


func reset(database *sql.DB) http.HandlerFunc {

	return func (w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/plain")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.SetFlags(log.Lshortfile)

		if r.Method == http.MethodPost {
			
			var id int
			var verify string
			var token string

			password := r.FormValue("currentPassword")
			newPassword := r.FormValue("newPassword")
			confirmPassword := r.FormValue("confirmPassword")

			email := teste[0]

			queryUpdatePassword, updatePasswordErro := database.Prepare("UPDATE usuarios SET password = ? WHERE email = ?")
			if updatePasswordErro != nil {
				log.Println("ERROR: ", updatePasswordErro)
			}

			query := "SELECT id, password FROM usuarios WHERE email = ?"
			queryPasswordError := database.QueryRow(query, email).Scan(&id, &verify)
			if queryPasswordError != nil {
				if queryPasswordError == sql.ErrNoRows {
					log.Println("PASSWORD NOT FOUND")
					return
				}
				log.Println("ERROR: ", queryPasswordError)
				return
			}

			getToken := "SELECT token FROM password_token WHERE user_id = ?"
			getTokenError := database.QueryRow(getToken, id).Scan(&token)
			if getTokenError != nil {
				if getTokenError == sql.ErrNoRows {
					log.Println("ID NOT FOUND")
					return
				}
				log.Println("ERROR: ", getTokenError)
				return
			}

			var idToken int
			var user_id int
			var expires_at string
			var used bool

			getInformation := "SELECT id, user_id, expires_at, used FROM password_token WHERE token = ?"
			getInformationError := database.QueryRow(getInformation, token).Scan(&idToken, &user_id, &expires_at, &used)
			if getInformationError != nil {
				if getInformationError == sql.ErrNoRows {
					log.Println("NOT FOUND")
					return
				}
				log.Println("ERROR: ", getInformationError)
				return
			}

			passwordHash := bcrypt.CompareHashAndPassword([]byte(verify), []byte(password))
			if passwordHash == nil && newPassword != "" && password != newPassword && utf8.RuneCountInString(newPassword) >= 8 && newPassword == confirmPassword {
				fmt.Println("VALID PASSWORD")

				hash, err := hashPassword(newPassword)
				if err != nil {
					log.Fatal("ERROR: ", err)
					return
				}

				_, errorPassword := queryUpdatePassword.Exec(hash, teste[0])
				if errorPassword != nil {
					log.Fatal("ERROR: ", errorPassword)
					return
				}

				_, tokenUsedError := database.Exec("UPDATE password_token SET used = TRUE WHERE user_id = ?", user_id)
				if tokenUsedError != nil {
					log.Fatal("ERROR: ", tokenUsedError)
					return 
				}

				context, cancel := context.WithCancel(context.Background())

				go func(context context.Context) {
					select {
					case <-time.After(20 * time.Second):
						teste = teste[:0]
						fmt.Println("")

					case <-ctx.Done():
						fmt.Println("")
					}
					
				}(context)
				
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("VALID PASSWORD"))

				cancel()

			} else if passwordHash != nil {
				fmt.Println("WRONG PASSWORD")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("INCORRECT PASSWORD"))

			} else if password == newPassword {
				fmt.Println("THE SAME PASSWORD")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("THE PASSWORD ARE THE SAME"))

			} else if utf8.RuneCountInString(newPassword) < 8 {
				fmt.Println("SHORT PASSWORD")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("SHORT PASSWORD"))

			} else if newPassword != confirmPassword {
				fmt.Println("PASSWORD CONFIRMATION IS WRONG")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("PASSWORD CONFIRMATION IS WRONG"))

			} else {
				log.Println("SOME ERROR")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("SOME ERROR"))
			}

		} else {
			http.Error(w, "METHOD NOT PERMITED", http.StatusMethodNotAllowed)
		}
	}
}


func generateTokens() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}


func generateLink(token string) string {
	return "" + token
}


func removeExpiredToken(database *sql.DB) {

	_, queryError := database.Exec("DELETE FROM password_token WHERE expires_at < NOW()")
	if queryError != nil {
		log.Println("ERROR: ", queryError)
		return
	}

	log.Println("TOKENS REMOVED")

}


func limitOfAttempts(database *sql.DB, email string) (bool, error) {

	var emailCount int

	err := database.QueryRow("SELECT COUNT(*) FROM attempts WHERE email = ? AND attempt_time > ?", email, time.Now().Add(-15*time.Minute)).Scan(&emailCount)
	if err != nil {
		return false, err
	}

	if emailCount >= 4 {
		return false, nil
	}

	return true, nil

}


func attemptLogs(database *sql.DB, email string) error {
	_, err := database.Exec("INSERT INTO attempts(email) VALUES(?)", email)
	return err
}


func startToRemoverExpiredTokens(database *sql.DB) {

	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			removeExpiredToken(database)
		}
	}()

}


func validTokenHandler(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/plain")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.SetFlags(log.Lshortfile)

		if r.Method == http.MethodPost {

			var id int
			var expires_at time.Time
			var used bool
	
			email := teste[0]
			err := database.QueryRow("SELECT id FROM usuarios WHERE email = ?", email).Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Println("EMAIL NOT FOUND")
					return
				}
				log.Println("ERROR: ", err)
				return
			}
	
			err = database.QueryRow("SELECT expires_at, used FROM password_token WHERE user_id = ?", id).Scan(&expires_at, &used)
			if err != nil {

				fmt.Println("EXPIRES AT: ", expires_at)
				fmt.Println("USED: ", used)
				if err == sql.ErrNoRows || used == true || time.Now().After(expires_at) {
					w.Write([]byte("INVALID TOKEN"))
					log.Println("INVALID TOKEN")
					return

				}
				log.Println("ERROR: ", err)
				return
	
			} else if err == nil {
				teste = teste[:0]
				w.Write([]byte("VALID TOKEN"))
				log.Println("VALID TOKEN")

			}
		}
	}
}


func main() {
	database := database()
	defer database.Close()

	http.HandleFunc("/sign", handler(db))
	http.HandleFunc("/login", handlerLogIn(db))

	http.HandleFunc("/change", handlerChangeName(db))
	http.HandleFunc("/email", handlerChangeEmail(db))

	http.HandleFunc("/logout", logoutHandler(db))
	http.HandleFunc("/delete", handlerDeleteAccount(db))

	http.HandleFunc("/reset", resetPassword(db))
	http.HandleFunc("/reset/password", reset(db))
	
	fmt.Println("SERVER OPEN WITH GOLANG")
	http.ListenAndServe("", nil)
}
