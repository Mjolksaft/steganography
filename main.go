package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"steganography/internal/auth"
	"steganography/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	mainMenu int = iota
	secondMenu
)

type CLIcommand struct {
	Name        string
	Description string
	Callback    func() error
}

var DB *sql.DB

func main() {
	godotenv.Load()
	// get database connection
	db, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))

	if err != nil {
		fmt.Println("error connecting to sql server:", err)
		os.Exit(0)
	}
	DB = db

	// create a server
	mux := http.NewServeMux()
	server := http.Server{Handler: mux, Addr: ":8080"}

	mux.Handle("/api/", http.StripPrefix("/api", http.FileServer(http.Dir("public"))))
	mux.HandleFunc("POST /api/passwords", createPasswordHandler)
	mux.HandleFunc("GET /api/passwords", getPasswordHandler)

	fmt.Println("now listening on port: 8080")
	server.ListenAndServe()

	// check login credentials

	// start CLIloop from mainmenu

	// start a api
}

func createPasswordHandler(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		Password    string `json:"password"`
		Application string `json:"application"`
	}

	// Create a decoder for the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close() // Ensure body is closed after reading

	// Decode the JSON body into the dataStruct
	var body dataStruct
	if err := decoder.Decode(&body); err != nil {
		fmt.Println(err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest) // Bad request if JSON is invalid
		return
	}
	// get the query
	dbQueries := database.New(DB)
	// complete the query
	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	newPassword, err := dbQueries.CreateOne(r.Context(), database.CreateOneParams{
		HashedPassword: sql.NullString{String: string(hashedPassword), Valid: true},
		Application:    sql.NullString{String: body.Application, Valid: true},
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error adding to database", http.StatusInternalServerError)
		return
	}

	// encode the new user
	encodedPassword, err := json.Marshal(newPassword)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error mashaling", http.StatusInternalServerError)
		return
	}

	// write result to user
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(encodedPassword))
}

func getPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("Hello there"))
}
