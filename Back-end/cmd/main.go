package main

import (
	"fmt"
	"net/http"

	"index/internal/database"
	"index/internal/handlers"
	"index/internal/service"
	"index/internal/repository"
)


func main() {
	database := database.Connect()
	defer database.Close()

	go service.StartToRemoverExpiredTokens()
	mux := http.NewServeMux()
	
	repositoryRegister := repository.NewRegisterStruct(database)
	serviceRegister := service.NewUserStruct(repositoryRegister)
	handlerRegister := handlers.NewRegisterHanlder(serviceRegister)

	repositoryLogin := repository.NewVerifyLoginStruct(database)
	serviceLogin := service.NewVerifyLogin(repositoryLogin)
	handlerLogin := handlers.NewLoginHandler(serviceLogin)

	repositoryChangeName := repository.NewChangeNameStruct(database)
	serviceChangeName := service.NewChangeName(repositoryChangeName)
	handlerChangeName := handlers.NewChangeNameHandler(serviceChangeName)

	repositoryChangeEmail := repository.NewChangeEmailStruct(database)
	serviceChangeEmail := service.NewChangeEmail(repositoryChangeEmail)
	handlerChangeEmail := handlers.NewChangeEmailHandler(serviceChangeEmail)

	repositoryRequest := repository.NewRequestStruct(database)
	serviceRequest := service.NewRequest(repositoryRequest)
	handlerRequest := handlers.NewRequestHandler(serviceRequest)

	repositoryResetPassword := repository.NewResetPasswordStruct(database)
	serviceResetPassword := service.NewResetPassword(repositoryResetPassword)
	handlerResetPassword := handlers.NewResetPasswordHandler(serviceResetPassword)

	repositoryDeleteAccount := repository.NewDeleteAccountStruct(database)
	serviceDeleteAccount := service.NewDeleteAccount(repositoryDeleteAccount)
	handlerDeleteAccount := handlers.NewDeleteAccountHandler(serviceDeleteAccount)

	repositoryValidToken := repository.NewValidTokenStruct(database)
	serviceValidToken := service.NewValidToken(repositoryValidToken)
	handlerValidToken := handlers.NewValidTokenHandler(serviceValidToken)


	mux.HandleFunc("/register", handlerRegister.RegisterHandler)
	mux.HandleFunc("/login", handlerLogin.HandlerLogin)

	mux.HandleFunc("/change", handlerChangeName.ChangeNameHandler)
	mux.HandleFunc("/email", handlerChangeEmail.ChangeEmailHandler)

	mux.HandleFunc("/delete", handlerDeleteAccount.DeleteAccountHandler)

	mux.HandleFunc("/reset", handlerRequest.RequestHandler)
	mux.HandleFunc("/reset/password", handlerResetPassword.ResetPasswordHandler)

	mux.HandleFunc("/valid", handlerValidToken.ValidTokenHandler)

	middleware := handlers.Recovery(mux)
	
	fmt.Println("SERVER OPEN WITH GOLANG")
	http.ListenAndServe("127.0.0.1:8000", middleware)
}
