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
	
	repo := &repository.MysqlStruct{}
	serviceSigUp := &service.User{Repository: repo}
	handler := &handlers.Handler{Service: serviceSigUp}

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

	http.HandleFunc("/sign", handler.NewSignUpHandler)
	http.HandleFunc("/login", handlerLogin.NewHandlerLogin)

	http.HandleFunc("/change", handlerChangeName.ChangeNameHandler)
	http.HandleFunc("/email", handlerChangeEmail.ChangeEmailHandler)

	http.HandleFunc("/logout", handlers.LogoutHandler(database))
	http.HandleFunc("/delete", handlers.HandlerDeleteAccount(database))

	http.HandleFunc("/reset", handlerRequest.RequestHandler)
	http.HandleFunc("/reset/password", handlerResetPassword.ResetPasswordHandler)
	http.HandleFunc("/reset/valid", handlers.ValidTokenHandler(database))
	
	fmt.Println("SERVER OPEN WITH GOLANG")
	http.ListenAndServe("127.0.0.1:8000", nil)
}
