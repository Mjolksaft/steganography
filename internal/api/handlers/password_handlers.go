package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"steganography/internal/database"
	"steganography/internal/middleware"

	"github.com/google/uuid"
)

type PasswordHandler struct {
	DB *sql.DB
}

// create password with password and application name
func (h PasswordHandler) CreatePassword(w http.ResponseWriter, r *http.Request) {
	// validate the cookie with middleware
	userID := r.Context().Value(middleware.KEY).(string)

	type dataStruct struct {
		Password    string `json:"password"`
		Application string `json:"application"`
	}

	// decode the body
	decoder := json.NewDecoder(r.Body)
	var data dataStruct
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "error decoding", http.StatusInternalServerError)
		return
	}

	//convert to uuid
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "error decoding", http.StatusInternalServerError)
		return
	}

	//make the request
	dbQueries := database.New(h.DB)
	arg := database.CreatePasswordParams{HashedPassword: data.Password, Application: data.Application, UserID: parsedID}
	if err = dbQueries.CreatePassword(r.Context(), arg); err != nil {
		http.Error(w, "error decoding", http.StatusInternalServerError)
		return
	}

	// response returns nothing
	w.WriteHeader(201)
}
