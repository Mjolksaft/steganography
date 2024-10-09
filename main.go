package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"steganography/internal/api/handlers"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type CLIcommand struct {
	Name        string
	Description string
	Callback    func() error
}

func main() {
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
	mux.HandleFunc("GET /api/passwords", passwordHandlers.GetPassword)

	userHandlers := handlers.UserHandler{DB: db}

	mux.HandleFunc("POST /api/login", userHandlers.Login)
	mux.HandleFunc("POST /api/users", userHandlers.CreateUser)
	mux.HandleFunc("GET /api/users", userHandlers.GetUser)

	fmt.Println("now listening on port: 8080")
	server.ListenAndServe()

	// check login credentials

	// start CLIloop from mainmenu

	// start a api
}
