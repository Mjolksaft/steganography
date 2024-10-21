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
	godotenv.Load()

	// get database connection and sessionManager
	sessionManager = auth.NewSessionManager()
	db, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))

	if err != nil {
		fmt.Println("error connecting to sql server:", err)
		os.Exit(0)
	}

	// create a server
	mux := http.NewServeMux()
	server := http.Server{Handler: mux, Addr: ":8080"}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("public"))))

	// users endpoints!!!!
	userHandlers := handlers.UserHandler{DB: db, SM: sessionManager}
	mux.HandleFunc("/api/login", userHandlers.Login)                                                          // No authentication required for login
	mux.HandleFunc("/api/users", userHandlers.CreateUser)                                                     // Create a user (public)
	mux.Handle("/api/users/{user_id}", middleware.ValidateSession(http.HandlerFunc(userHandlers.UpdateUser))) // Update own user or as admin (authenticated)
	mux.Handle("/api/users/{user_id}", middleware.ValidateSession(http.HandlerFunc(userHandlers.DeleteUser))) // Delete own user or as admin (authenticated)
	mux.Handle("/api/users/{user_id}", middleware.ValidateSession(http.HandlerFunc(userHandlers.GetUser)))    // Get user data (authenticated)

	// password endpoints!!!!
	passwordHandlers := handlers.PasswordHandler{DB: db, SM: sessionManager}
	mux.Handle("/api/passwords", middleware.ValidateSession(http.HandlerFunc(passwordHandlers.CreatePassword)))               // Create password entry (authenticated)
	mux.Handle("/api/passwords/{password_id}", middleware.ValidateSession(http.HandlerFunc(passwordHandlers.UpdatePassword))) // Update password entry (authenticated)
	mux.Handle("/api/passwords/{password_id}", middleware.ValidateSession(http.HandlerFunc(passwordHandlers.DeletePassword))) // Delete password entry (authenticated)
	mux.Handle("/api/passwords", middleware.ValidateSession(http.HandlerFunc(passwordHandlers.GetPassword)))                  // Get password entry (authenticated)

	// admin endpoints!!!!
	adminHandlers := handlers.AdminHandler{DB: db, SM: sessionManager}
	mux.Handle("/admin/users", middleware.ValidateAdmin(http.HandlerFunc(adminHandlers.CreateUser)))                       // Create user (admin only)
	mux.Handle("/admin/users/{user_id}", middleware.ValidateAdmin(http.HandlerFunc(adminHandlers.UpdateUser)))             // Update any user (admin only)
	mux.Handle("/admin/users/{user_id}", middleware.ValidateAdmin(http.HandlerFunc(adminHandlers.DeleteUser)))             // Delete any user (admin only)
	mux.Handle("/admin/users", middleware.ValidateAdmin(http.HandlerFunc(adminHandlers.GetUsers)))                         // Get all users (admin only)
	mux.Handle("/admin/passwords/{password_id}", middleware.ValidateAdmin(http.HandlerFunc(adminHandlers.UpdatePassword))) // Update any user's password (admin only)
	mux.Handle("/admin/passwords/{password_id}", middleware.ValidateAdmin(http.HandlerFunc(adminHandlers.DeletePassword))) // Delete any user's password (admin only)
	mux.Handle("/admin/passwords", middleware.ValidateAdmin(http.HandlerFunc(adminHandlers.GetPasswords)))                 // Get all passwords (admin only)

	fmt.Println("now listening on port: 8080")
	server.ListenAndServe()
}
