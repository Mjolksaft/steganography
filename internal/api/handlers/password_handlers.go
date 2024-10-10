package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"steganography/internal/auth"
	"steganography/internal/database"
)

type PasswordHandler struct {
	DB *sql.DB
}

func (h PasswordHandler) CreatePassword(w http.ResponseWriter, r *http.Request) {
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
	dbQueries := database.New(h.DB)
	// complete the query
	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = dbQueries.CreatePassword(r.Context(), database.CreatePasswordParams{
		HashedPassword: sql.NullString{String: string(hashedPassword), Valid: true},
		Application:    sql.NullString{String: body.Application, Valid: true},
	})

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error adding to database", http.StatusInternalServerError)
		return
	}

	// write result to user
	w.WriteHeader(200)
}

func (h PasswordHandler) GetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("Hello there"))
}
