package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"steganography/internal/auth"
	"steganography/internal/database"
)

type UserHandler struct {
	DB *sql.DB
}

func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// check the body
	type password struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	//decode the body
	decoder := json.NewDecoder(r.Body)

	var data password
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "error decoding body", http.StatusInternalServerError)
		return
	}

	// make the query
	dbQueries := database.New(h.DB)
	user, err := dbQueries.GetUser(r.Context(), sql.NullString{String: data.Username, Valid: true})
	if err != nil {
		http.Error(w, "wrong username or password", http.StatusBadRequest)
		return
	}

	// check the password
	err = auth.CheckPassword(user.HashedPassword.String, data.Password)
	if err != nil {
		http.Error(w, "password doesnot match", http.StatusBadRequest)
		return
	}

	// send responese
	w.Header().Add("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Username: %s, Password %s", data.Username, data.Password)))
}

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// decode the body
	type password struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var data password
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "error decoding body", http.StatusBadRequest)
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
	}

	// add to database
	dbQueries := database.New(h.DB)
	user, err := dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Username:       sql.NullString{String: data.Username, Valid: true},
		HashedPassword: sql.NullString{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		http.Error(w, "error creating user", http.StatusBadRequest)
	}

	// encode user
	encodedData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error encoding data", http.StatusInternalServerError)
	}

	// write response
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	w.Write(encodedData)
}

func (h UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

}
