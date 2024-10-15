package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"steganography/internal/api/handlers"
	"steganography/internal/auth"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var sessionManager *auth.SessionManager

func main() {
	//inject the sessionManager to the handlers that need it
	sessionManager = auth.NewSessionManager()
	handlers.InitSessionManager(sessionManager)

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
	passwordHandlers := handlers.PasswordHandler{DB: db}

	mux.HandleFunc("POST /api/passwords", passwordHandlers.CreatePassword)

	userHandlers := handlers.UserHandler{DB: db}

	mux.HandleFunc("POST /api/login", userHandlers.Login)
	mux.HandleFunc("POST /api/users", userHandlers.CreateUser)
	mux.HandleFunc("GET /api/users", userHandlers.GetUser)

	fmt.Println("now listening on port: 8080")
	server.ListenAndServe()
}
