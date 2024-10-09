package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"steganography/internal/auth"
	"steganography/internal/database"
	"time"
)

type UserHandler struct {
	DB *sql.DB
}

type User struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	HashedPassword string    `json:"hashed_password"`
	Username       string    `json:"username"`
	IsAdmin        bool      `json:"is_admin"`
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Define a struct to capture the request body
	type credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decode the body
	var data credentials
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close() // Ensure the body is closed after decoding
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "error decoding body", http.StatusBadRequest)
		return
	}

	// Query the user by username
	dbQueries := database.New(h.DB)
	user, err := dbQueries.GetUser(r.Context(), sql.NullString{String: data.Username, Valid: true}) // No need for sql.NullString
	if err != nil {
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	// Check the password
	if err := auth.CheckPassword(user.HashedPassword.String, data.Password); err != nil {
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		return
	}

	// Map the data to the User struct
	loggedUser := User{
		ID:             user.ID.String(),
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
		Username:       user.Username.String,
		HashedPassword: user.HashedPassword.String, // It's not necessary to send this back
	}

	//encode the user
	encoded, err := json.Marshal(loggedUser)
	if err != nil {
		http.Error(w, "could not marshal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK) // Use a constant for the status code
	w.Write(encoded)
}

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// decode the body
	type password struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var data password
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "error decoding body", http.StatusBadRequest)
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	// add to database
	dbQueries := database.New(h.DB)
	user, err := dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Username:       sql.NullString{String: data.Username, Valid: true},
		HashedPassword: sql.NullString{String: string(hashedPassword), Valid: true},
	})

	if err != nil {
		http.Error(w, "error creating user", http.StatusBadRequest)
		return
	}

	// encode user
	encodedData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error encoding data", http.StatusInternalServerError)
		return
	}

	// write response
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	w.Write(encodedData)
}

func (h UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

}
