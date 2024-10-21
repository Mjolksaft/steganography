package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
	arg := database.CreatePasswordParams{HashedPassword: data.Password, ApplicationName: data.Application, UserID: parsedID}
	if err = dbQueries.CreatePassword(r.Context(), arg); err != nil {
		http.Error(w, "error decoding", http.StatusInternalServerError)
		return
	}

	// response returns nothing
	w.WriteHeader(201)
}

func (h PasswordHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {

}

func (h PasswordHandler) DeletePassword(w http.ResponseWriter, r *http.Request) {

}

// Retrieves a password entry for a specific application for the authenticated user. If no application_name is provided, return all passwords for the user.
func (h PasswordHandler) GetPassword(w http.ResponseWriter, r *http.Request) {
	// read the query param and context key
	application := r.FormValue("application_name")
	userID, ok := r.Context().Value(middleware.KEY).(string)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	//convert to uuid
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "error decoding", http.StatusInternalServerError)
		return
	}

	if application != "" {
		// return specific password
		dbQueries := database.New(h.DB)
		data, err := dbQueries.GetPassword(r.Context(), database.GetPasswordParams{ApplicationName: application, UserID: parsedID})
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error quering database", http.StatusInternalServerError)
			return
		}

		//encode data
		dataByte, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error encoding data", http.StatusInternalServerError)
			return
		}

		// send data
		w.Header().Set("content-type", "application/json; charset=utf8")
		w.WriteHeader(200)
		w.Write(dataByte)

	} else {
		// return passwords
		dbQueries := database.New(h.DB)
		data, err := dbQueries.GetPasswords(r.Context(), parsedID)
		fmt.Println(data)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error quering database", http.StatusInternalServerError)
			return
		}

		//encode data
		dataByte, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error encoding data", http.StatusInternalServerError)
			return
		}

		// send data
		w.Header().Set("content-type", "application/json; charset=utf8")
		w.WriteHeader(200)
		w.Write(dataByte)
	}
}
