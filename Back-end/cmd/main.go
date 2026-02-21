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

	repositoryLogin := &repository.VerifyLoginStruct{}
	serviceLogin := &service.VerifyLogin{Repository: repositoryLogin}
	handlerLogin := &handlers.LoginHandler{Service: serviceLogin}

	repositoryChangeName := &repository.ChangeNameStruct{}
	serviceChangeName := &service.ChangeName{Repository: repositoryChangeName}	
	handlerChangeName := &handlers.ChangeNameHandler{Service: serviceChangeName}

	repositoryChangeEmail := &repository.ChangeEmailStruct{}
	serviceChangeEmail := &service.ChangeEmail{Repository: repositoryChangeEmail}
	handlerChangeEmail := &handlers.ChangeEmailHandler{Service: serviceChangeEmail}

	repositoryRequest := &repository.RequestStruct{}
	serviceRequest := &service.Request{Repository: repositoryRequest}
	handlerRequest := &handlers.RequestHandler{Service: serviceRequest}

	repositoryResetPassword := &repository.ResetPasswordStruct{}
	serviceResetPassword := &service.ResetPassword{Repository: repositoryResetPassword}
	handlerResetPassword := &handlers.ResetPasswordHandler{Service: serviceResetPassword}

	repositoryDeleteAccount := &repository.DeleteAccountStruct{}
	serviceDeleteAccount := &service.DeleteAccount{Repository: repositoryDeleteAccount}
	handlerDeleteAccount := &handlers.DeleteAccountHandler{Service: serviceDeleteAccount}

	mux.HandleFunc("/register", handlerRegister.RegisterHandler)
	mux.HandleFunc("/login", handlerLogin.HandlerLogin)

	mux.HandleFunc("/change", handlerChangeName.ChangeNameHandler)
	mux.HandleFunc("/email", handlerChangeEmail.ChangeEmailHandler)

	mux.HandleFunc("/delete", handlerDeleteAccount.DeleteAccountHandler)

	mux.HandleFunc("/reset", handlerRequest.RequestHandler)
	mux.HandleFunc("/reset/password", handlerResetPassword.ResetPasswordHandler)

	middleware := handlers.Recovery(mux)
	
	fmt.Println("SERVER OPEN WITH GOLANG")
	http.ListenAndServe("127.0.0.1:8000", middleware)
}
