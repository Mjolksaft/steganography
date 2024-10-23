package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"steganography/internal/auth"
	"steganography/internal/database"
	"steganography/internal/middleware"
	"time"
)

type UserHandler struct {
	DB *sql.DB
	SM *auth.SessionManager
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
		slog.Error("Failed to decode login request body", slog.String("err", err.Error()))
		return
	}

	slog.Info("Login attempt", slog.String("username", data.Username))

	// Query the user by username
	dbQueries := database.New(h.DB)
	user, err := dbQueries.GetUser(r.Context(), data.Username)
	if err != nil {
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		slog.Error("Failed login: user not found or incorrect password", slog.String("username", data.Username), slog.String("err", err.Error()))
		return
	}

	// Check the password
	if err := auth.CheckPassword(user.HashedPassword, data.Password); err != nil {
		http.Error(w, "incorrect username or password", http.StatusUnauthorized)
		slog.Error("Failed login: incorrect password", slog.String("username", data.Username), slog.String("err", err.Error()))
		return
	}

	slog.Info("Login successful", slog.String("username", data.Username))

	// Map the data to the User struct
	loggedUser := User{
		ID:             user.ID.String(),
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
	}

	// Encode the user
	encoded, err := json.Marshal(loggedUser)
	if err != nil {
		http.Error(w, "could not marshal", http.StatusInternalServerError)
		slog.Error("Failed to encode logged-in user data", slog.String("username", data.Username), slog.String("err", err.Error()))
		return
	}

	// Create session token
	expiredAt := time.Now().Add(time.Second * 120)
	sessionToken := h.SM.CreateSession(loggedUser.ID, expiredAt)

	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiredAt,
	})

	slog.Info("Session created", slog.String("username", data.Username), slog.String("sessionToken", sessionToken))

	// Send the response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Define the struct for user credentials
	type password struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decode the request body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var data password
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "error decoding body", http.StatusBadRequest)
		slog.Error("Failed to decode CreateUser request body", slog.String("err", err.Error()))
		return
	}

	slog.Info("CreateUser attempt", slog.String("username", data.Username))

	// Hash the password
	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		slog.Error("Failed to hash password", slog.String("username", data.Username), slog.String("err", err.Error()))
		return
	}

	// Add user to the database
	dbQueries := database.New(h.DB)
	user, err := dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Username:       data.Username,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "error creating user", http.StatusBadRequest)
		slog.Error("Failed to create user in the database", slog.String("username", data.Username), slog.String("err", err.Error()))
		return
	}

	slog.Info("User created successfully", slog.String("username", data.Username))

	// Encode user data for response
	encodedData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error encoding data", http.StatusInternalServerError)
		slog.Error("Failed to encode user data for response", slog.String("username", data.Username), slog.String("err", err.Error()))
		return
	}

	// Write response
	w.Header().Add("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(encodedData)
}

func (h UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve userID from the context
	userID, ok := r.Context().Value(middleware.KEY).(string)
	fmt.Println(userID, ok)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	// Use userID for logic (e.g., fetching user details from the database)
	fmt.Fprintf(w, "User ID from context: %s", userID)
}

func (h UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (h UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}
