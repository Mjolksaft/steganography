package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"steganography/internal/auth"
	"steganography/internal/database"
	"steganography/internal/middleware"

	"github.com/google/uuid"
)

type PasswordHandler struct {
	DB *sql.DB
	SM *auth.SessionManager
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
		slog.Error("Failed to decode request body", slog.String("userID", userID), slog.String("err", err.Error()))
		return
	}

	// URL-decode the application name
	decodedAppName, err := url.QueryUnescape(data.Application)
	if err != nil {
		http.Error(w, "Invalid application name", http.StatusBadRequest)
		slog.Error("Invalid application name", slog.String("application", data.Application), slog.String("userID", userID), slog.String("err", err.Error()))
		return
	}

	// Convert to uuid
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "error decoding user ID", http.StatusInternalServerError)
		slog.Error("Failed to parse user ID", slog.String("userID", userID), slog.String("err", err.Error()))
		return
	}

	slog.Info("CreatePassword request received", slog.String("application", decodedAppName), slog.String("userID", parsedID.String()))

	// Create the password in the database
	dbQueries := database.New(h.DB)
	arg := database.CreatePasswordParams{
		HashedPassword:  data.Password,
		ApplicationName: decodedAppName,
		UserID:          parsedID,
	}

	if err = dbQueries.CreatePassword(r.Context(), arg); err != nil {
		http.Error(w, "error saving to database", http.StatusInternalServerError)
		slog.Error("Failed to create password in database", slog.String("application", decodedAppName), slog.String("userID", parsedID.String()), slog.String("err", err.Error()))
		return
	}

	slog.Info("Password successfully created", slog.String("application", decodedAppName), slog.String("userID", parsedID.String()))

	// Send success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Password created successfully"))
}

func (h PasswordHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {

}

func (h PasswordHandler) DeletePassword(w http.ResponseWriter, r *http.Request) {
	// Extract password_id from the URL variables
	passwordId := r.PathValue("password_id")

	// get userId from cotnext
	userID, ok := r.Context().Value(middleware.KEY).(string)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		slog.Error("User not found in context", slog.String("contextKey", string(middleware.KEY)))
		return
	}
	fmt.Println(passwordId, userID)
	//convert to uuid
	parsedPasswordId, err := uuid.Parse(passwordId)
	if err != nil {
		http.Error(w, "parsing to uuid failed", http.StatusInternalServerError)
		slog.Error("parsing to uuid failed", slog.String("contextKey", string(middleware.KEY)))
	}
	parsedUserId, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "parsing to uuid failed", http.StatusInternalServerError)
		slog.Error("parsing to uuid failed", slog.String("contextKey", string(middleware.KEY)))
	}

	// create query with user id and password id
	dbQueries := database.New(h.DB)
	err = dbQueries.DeletePassword(r.Context(), database.DeletePasswordParams{ID: parsedPasswordId, UserID: parsedUserId})
	if err != nil {
		http.Error(w, "Query failed in db", http.StatusInternalServerError)
		slog.Error("Query failed in db", slog.String("contextKey", string(middleware.KEY)))
	}

	slog.Error("password deleted from db", slog.String("contextKey", string(middleware.KEY)))
	//return result
	w.WriteHeader(200)
}

// Retrieves a password entry for a specific application for the authenticated user. If no application_name is provided, return all passwords for the user.
func (h PasswordHandler) GetPassword(w http.ResponseWriter, r *http.Request) {
	// read the query param and context key
	application := r.FormValue("application_name")
	userID, ok := r.Context().Value(middleware.KEY).(string)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		slog.Error("User not found in context", slog.String("contextKey", string(middleware.KEY)))
		return
	}

	// URL-decode the application name
	decodedAppName, err := url.QueryUnescape(application)
	if err != nil {
		http.Error(w, "Invalid application name", http.StatusBadRequest)
		slog.Error("Invalid application name", slog.String("application", application), slog.String("err", err.Error()))
		return
	}

	// Convert to uuid
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Error decoding user ID", http.StatusInternalServerError)
		slog.Error("Failed to parse user ID", slog.String("userID", userID), slog.String("err", err.Error()))
		return
	}

	slog.Info("GetPassword called", slog.String("application", decodedAppName), slog.String("userID", parsedID.String()))

	dbQueries := database.New(h.DB)

	if application != "" {
		// Return specific password
		data, err := dbQueries.GetPassword(r.Context(), database.GetPasswordParams{ApplicationName: decodedAppName, UserID: parsedID})
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			slog.Error("Database query failed for GetPassword", slog.String("application", decodedAppName), slog.String("userID", parsedID.String()), slog.String("err", err.Error()))
			return
		}

		// Encode data
		dataByte, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error encoding data", http.StatusInternalServerError)
			slog.Error("Failed to encode password data", slog.String("application", decodedAppName), slog.String("userID", parsedID.String()), slog.String("err", err.Error()))
			return
		}

		// Send response
		w.Header().Set("content-type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write(dataByte)

		// Log successful response
		slog.Info("Password returned successfully", slog.String("application", decodedAppName), slog.String("userID", parsedID.String()))

	} else {
		// Return all passwords
		data, err := dbQueries.GetPasswords(r.Context(), parsedID)
		if err != nil {
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			slog.Error("Database query failed for GetPasswords", slog.String("userID", parsedID.String()), slog.String("err", err.Error()))
			return
		}

		// Encode data
		dataByte, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error encoding data", http.StatusInternalServerError)
			slog.Error("Failed to encode passwords data", slog.String("userID", parsedID.String()), slog.String("err", err.Error()))
			return
		}

		// Send response
		w.Header().Set("content-type", "application/json; charset=utf8")
		w.WriteHeader(http.StatusOK)
		w.Write(dataByte)

		slog.Info("Passwords returned successfully", slog.String("userID", parsedID.String()))
	}
}
