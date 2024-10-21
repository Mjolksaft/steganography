package handlers

import (
	"database/sql"
	"net/http"
	"steganography/internal/auth"
)

type AdminHandler struct {
	DB *sql.DB
	SM *auth.SessionManager
}

func (h AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

}

func (h AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (h AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (h AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

}

func (h AdminHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {

}

func (h AdminHandler) DeletePassword(w http.ResponseWriter, r *http.Request) {

}

func (h AdminHandler) GetPasswords(w http.ResponseWriter, r *http.Request) {

}
