package main

import (
	"context"
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
	"github.com/redis/go-redis/v9"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Warning: .env file was not found")
	}

	jwtKey := os.Getenv("JWT_KEY")
	redisAddr := os.Getenv("REDIS_ADDR")

	limiter, err := security.NewRedisLimiter(redisAddr)
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	db := database.Connect()
	defer db.Close()
	database.RunMigrations(db)

	// go service.StartToRemoverExpiredTokens(db)
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
	sec := security.Security()

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr, DB: 0})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}

	resetStore := security.NewResetPassword(rdb)
	serviceRequest := service.NewService(
		userRepo,
		sec,
		resetStore,
		limiter,
	)
	handlerRequest := handlers.NewRequestHandler(serviceRequest)
	handlerValidToken := handlers.NewValidTokenHandler(serviceRequest)

	repositoryResetPassword := repository.NewResetPasswordStruct(db)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}

	securityReset := security.NewResetPassword(rdb)

	serviceResetPassword := service.NewResetPassword(repositoryResetPassword, securityReset)
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
	mux.Handle("/reset/password", middleware.AuthMiddleware(jwtKey)(http.HandlerFunc(handlerResetPassword.ResetPasswordHandler)))

	mux.HandleFunc("/valid", handlerValidToken.ValidToken)

	handlersWithRecovery := middleware.Recovery(mux)
	middleware := middleware.CorsMiddleware(handlersWithRecovery)

	log.Println("SERVER OPEN WITH GOLANG")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", middleware))
}
