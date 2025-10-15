// backend/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Arga-12/SimpleNotesSharingApp/app/backend/handlers"
	"github.com/Arga-12/SimpleNotesSharingApp/app/backend/middlewares"
)

func main() {
	// Load server configuration
	serverPort := getenvLocal("PORT", "8080")
	db := InitDB()
	defer db.Close()

	// Set database for logging middleware
	middlewares.SetLogDB(db)

	// Initialize handlers
	authHandler := &handlers.AuthHandler{DB: db}
	notesHandler := &handlers.NotesHandler{DB: db}

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes
	mux.Handle("/api/register", middlewares.Logging(http.HandlerFunc(authHandler.HandleRegister)))
	mux.Handle("/api/login", middlewares.Logging(http.HandlerFunc(authHandler.HandleLogin)))
	mux.Handle("/api/me", middlewares.Logging(http.HandlerFunc(authHandler.HandleMe)))
	mux.Handle("/api/logout", middlewares.Logging(http.HandlerFunc(authHandler.HandleLogout)))

	// Notes routes
	mux.Handle("/api/notes", middlewares.Logging(http.HandlerFunc(notesHandler.HandleNotes)))
	mux.Handle("/api/notes/", middlewares.Logging(http.HandlerFunc(notesHandler.HandleNoteByID)))

	addr := ":" + serverPort
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, middlewares.AllowLocalhostCookies(mux)); err != nil {
		log.Fatalf("server: %v", err)
	}
}
func getenvLocal(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
