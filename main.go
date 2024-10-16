package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"steganography/internal/api/handlers"
	"steganography/internal/auth"
	"steganography/internal/middleware"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var sessionManager *auth.SessionManager

func main() {
	//inject the sessionManager to the handlers that need it
	sessionManager = auth.NewSessionManager()
	handlers.InitSessionManager(sessionManager)
	middleware.InitSessionManager(sessionManager)

	godotenv.Load()

	// get database connection
	db, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))

	if err != nil {
		fmt.Println("error connecting to sql server:", err)
		os.Exit(0)
	}

	// create a server
	mux := http.NewServeMux()
	server := http.Server{Handler: mux, Addr: ":8080"}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("public"))))

	//password endpoits!!!!
	passwordHandlers := handlers.PasswordHandler{DB: db}
	mux.Handle("POST /api/passwords", middleware.ValidateSession(http.HandlerFunc(passwordHandlers.CreatePassword))) // Apply middleware
	// mux.HandleFunc("POST /api/passwords", passwordHandlers.CreatePassword)

	//users endpoits!!!!
	userHandlers := handlers.UserHandler{DB: db}
	mux.HandleFunc("POST /api/login", userHandlers.Login)
	mux.HandleFunc("POST /api/users", userHandlers.CreateUser)
	mux.Handle("GET /api/users", middleware.ValidateSession(http.HandlerFunc(userHandlers.GetUser))) // Apply middleware

	fmt.Println("now listening on port: 8080")
	server.ListenAndServe()
}
