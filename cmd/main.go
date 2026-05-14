package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ShieldAuth-API/internal/database"
	"ShieldAuth-API/internal/handlers"
	"ShieldAuth-API/internal/middleware"
	"ShieldAuth-API/internal/repository"
	"ShieldAuth-API/internal/security"
	"ShieldAuth-API/internal/service"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file was not found")
	}
	jwtKey := os.Getenv("JWT_KEY")

	limiter, err := security.NewRedisLimiter("")
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	db := database.Connect()
	defer db.Close()
	database.RunMigrations(db)

	go service.StartToRemoverExpiredTokens(db)
	mux := http.NewServeMux()

	repositoryRegister := repository.NewRegisterStruct(db)
	serviceRegister := service.NewUserStruct(repositoryRegister)
	handlerRegister := handlers.NewRegisterHanlder(serviceRegister)

	repositoryLogin := repository.NewVerifyLoginStruct(db)
	serviceLogin := service.NewVerifyLogin(repositoryLogin)
	handlerLogin := handlers.NewLoginHandler(serviceLogin, limiter)

	repositoryChangeName := repository.NewChangeNameStruct(db)
	serviceChangeName := service.NewChangeName(repositoryChangeName)
	handlerChangeName := handlers.NewChangeNameHandler(serviceChangeName)

	repositoryChangeEmail := repository.NewChangeEmailStruct(db)
	serviceChangeEmail := service.NewChangeEmail(repositoryChangeEmail)
	handlerChangeEmail := handlers.NewChangeEmailHandler(serviceChangeEmail)

	userRepo := repository.NewRequestStruct(db)
	tokenRepo := repository.NewValidTokenStruct(db)
	security := security.Security()
	serviceRequest := service.NewService(
		userRepo,
		tokenRepo,
		security,
	)
	handlerRequest := handlers.NewRequestHandler(serviceRequest)
	handlerValidToken := handlers.NewValidTokenHandler(serviceRequest)

	repositoryResetPassword := repository.NewResetPasswordStruct(db)
	serviceResetPassword := service.NewResetPassword(repositoryResetPassword)
	handlerResetPassword := handlers.NewResetPasswordHandler(serviceResetPassword)

	repositoryDeleteAccount := repository.NewDeleteAccountStruct(db)
	serviceDeleteAccount := service.NewDeleteAccount(repositoryDeleteAccount)
	handlerDeleteAccount := handlers.NewDeleteAccountHandler(serviceDeleteAccount)

	mux.HandleFunc("/register", handlerRegister.RegisterHandler)	
	mux.HandleFunc("/login", handlerLogin.HandlerLogin)

	auth := middleware.AuthMiddleware(jwtKey)
	mux.Handle("/change", auth(http.HandlerFunc(handlerChangeName.ChangeNameHandler)))

	handler := http.HandlerFunc(handlerChangeEmail.ChangeEmailHandler)
	handlerWithMiddleware := middleware.Chain(
		handler,
		middleware.AuthMiddleware(jwtKey),
		middleware.RateLimitMiddleware(limiter, "change-email-attempt", 3, 24*time.Hour),
	)
	mux.HandleFunc("/email", handlerWithMiddleware.ServeHTTP)
	// mux.HandleFunc("/email", handlerChangeEmail.ChangeEmailHandler)

	mux.HandleFunc("/delete", handlerDeleteAccount.DeleteAccountHandler)

	mux.HandleFunc("/reset", handlerRequest.RequestReset)
	mux.HandleFunc("/reset/password", handlerResetPassword.ResetPasswordHandler)

	mux.HandleFunc("/valid", handlerValidToken.ValidToken)

	handlersWithRecovery := middleware.Recovery(mux)
	middleware := middleware.CorsMiddleware(handlersWithRecovery)

	fmt.Println("SERVER OPEN WITH GOLANG")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", middleware))
}
