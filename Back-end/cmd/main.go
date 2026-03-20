package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"index/Back-end/internal/database"
	"index/Back-end/internal/handlers"
	"index/Back-end/internal/repository"
	"index/Back-end/internal/service"

	"github.com/joho/godotenv"
)


func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: file .env was not found")
	}
	jwtKey := os.Getenv("JWT_KEY")

	db := database.Connect()
	defer db.Close()
	// database.RunMigrations(db)

	go service.StartToRemoverExpiredTokens(db)
	mux := http.NewServeMux()
	
	repositoryRegister := repository.NewRegisterStruct(db)
	serviceRegister := service.NewUserStruct(repositoryRegister)
	handlerRegister := handlers.NewRegisterHanlder(serviceRegister)

	repositoryLogin := repository.NewVerifyLoginStruct(db)
	serviceLogin := service.NewVerifyLogin(repositoryLogin)
	handlerLogin := handlers.NewLoginHandler(serviceLogin)

	repositoryChangeName := repository.NewChangeNameStruct(db)
	serviceChangeName := service.NewChangeName(repositoryChangeName)
	handlerChangeName := handlers.NewChangeNameHandler(serviceChangeName)

	repositoryChangeEmail := repository.NewChangeEmailStruct(db)
	serviceChangeEmail := service.NewChangeEmail(repositoryChangeEmail)
	handlerChangeEmail := handlers.NewChangeEmailHandler(serviceChangeEmail)

	repositoryRequest := repository.NewRequestStruct(db)
	serviceRequest := service.NewRequest(repositoryRequest)
	handlerRequest := handlers.NewRequestHandler(serviceRequest)

	repositoryResetPassword := repository.NewResetPasswordStruct(db)
	serviceResetPassword := service.NewResetPassword(repositoryResetPassword)
	handlerResetPassword := handlers.NewResetPasswordHandler(serviceResetPassword)

	repositoryDeleteAccount := repository.NewDeleteAccountStruct(db)
	serviceDeleteAccount := service.NewDeleteAccount(repositoryDeleteAccount)
	handlerDeleteAccount := handlers.NewDeleteAccountHandler(serviceDeleteAccount)

	repositoryValidToken := repository.NewValidTokenStruct(db)
	serviceValidToken := service.NewValidToken(repositoryValidToken)
	handlerValidToken := handlers.NewValidTokenHandler(serviceValidToken)


	mux.HandleFunc("/register", handlerRegister.RegisterHandler)
	mux.HandleFunc("/login", handlerLogin.HandlerLogin)

	mux.HandleFunc("/change", handlerChangeName.ChangeNameHandler)

	getEnvMiddleware := handlers.AuthMiddleware(jwtKey)
	mux.Handle("/email", getEnvMiddleware(http.HandlerFunc(handlerChangeEmail.ChangeEmailHandler)))
	
	mux.HandleFunc("/delete", handlerDeleteAccount.DeleteAccountHandler)

	mux.HandleFunc("/reset", handlerRequest.RequestHandler)
	mux.HandleFunc("/reset/password", handlerResetPassword.ResetPasswordHandler)

	mux.HandleFunc("/valid", handlerValidToken.ValidTokenHandler)

	testHandler := handlers.Recovery(mux)
	middleware := handlers.CorsMiddleware(testHandler)
	
	fmt.Println("SERVER OPEN WITH GOLANG")
	http.ListenAndServe("127.0.0.1:8000", middleware)
}
