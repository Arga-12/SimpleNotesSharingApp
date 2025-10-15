package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Arga-12/SimpleNotesSharingApp/app/backend/middlewares"
	"github.com/Arga-12/SimpleNotesSharingApp/app/backend/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		jsonError(w, "username & password required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "server error", http.StatusInternalServerError)
		return
	}

	var id int
	var createdAt time.Time
	q := `INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id, created_at`
	err = h.DB.QueryRow(q, req.Username, string(hash), req.Email).Scan(&id, &createdAt)
	if err != nil {
		jsonError(w, "could not create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, err := middlewares.CreateJWT(id, req.Username)
	if err != nil {
		jsonError(w, "failed to create jwt", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	u := models.User{ID: id, Username: req.Username, Email: req.Email, CreatedAt: createdAt}
	jsonResponse(w, u, http.StatusCreated)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid json", http.StatusBadRequest)
		return
	}

	var id int
	var hashed string
	var username string
	q := `SELECT id, username, password FROM users WHERE username=$1`
	err := h.DB.QueryRow(q, req.Username).Scan(&id, &username, &hashed)
	if err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(req.Password)); err != nil {
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	tokenString, err := middlewares.CreateJWT(id, username)
	if err != nil {
		jsonError(w, "server error", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	jsonResponse(w, map[string]string{"message": "logged in"}, http.StatusOK)
}

func (h *AuthHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := middlewares.ParseJWT(c.Value)
	if err != nil {
		jsonError(w, "invalid token", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		jsonError(w, "invalid token claims", http.StatusUnauthorized)
		return
	}

	uid, _ := strconv.Atoi(sub)
	var u models.User
	err = h.DB.QueryRow(`SELECT id, username, email, created_at FROM users WHERE id=$1`, uid).
		Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, u, http.StatusOK)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, c)
	jsonResponse(w, map[string]string{"message": "logged out"}, http.StatusOK)
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
