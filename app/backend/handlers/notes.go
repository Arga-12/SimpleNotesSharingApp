package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Arga-12/SimpleNotesSharingApp/app/backend/middlewares"
	"github.com/Arga-12/SimpleNotesSharingApp/app/backend/models"
)

type NotesHandler struct {
	DB *sql.DB
}

func (h *NotesHandler) HandleNotes(w http.ResponseWriter, r *http.Request) {
	// All notes endpoints require JWT
	userID, err := middlewares.GetUserIDFromCookie(r)
	if err != nil {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// GET /notes → Return ALL notes from all users
		rows, err := h.DB.Query(`
			SELECT n.id, n.owner_id, u.username, n.title, n.content, n.shared, n.favorite, n.updated_at
			FROM notes n
			JOIN users u ON u.id = n.owner_id
			ORDER BY n.updated_at DESC`)
		if err != nil {
			jsonError(w, "db error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notes []models.Note
		for rows.Next() {
			var n models.Note
			if err := rows.Scan(&n.ID, &n.OwnerID, &n.OwnerUsername, &n.Title, &n.Content, &n.Shared, &n.Favorite, &n.Updated); err == nil {
				notes = append(notes, n)
			}
		}
		if len(notes) == 0 {
			jsonResponse(w, []models.Note{}, http.StatusOK)
			return
		}
		jsonResponse(w, notes, http.StatusOK)

	case http.MethodPost:
		var req models.Note
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "invalid json", http.StatusBadRequest)
			return
		}
		if req.Title == "" {
			req.Title = "Untitled Note"
		}

		var id int
		q := `INSERT INTO notes (owner_id, title, content, shared, favorite, updated_at)
			  VALUES ($1, $2, $3, $4, $5, now()) RETURNING id`
		err := h.DB.QueryRow(q, userID, req.Title, req.Content, req.Shared, req.Favorite).Scan(&id)
		if err != nil {
			jsonError(w, "db insert error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		req.ID = id
		req.OwnerID = userID
		req.Updated = time.Now()
		_ = h.DB.QueryRow(`SELECT username FROM users WHERE id=$1`, userID).Scan(&req.OwnerUsername)
		jsonResponse(w, req, http.StatusCreated)

	default:
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *NotesHandler) HandleNoteByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserIDFromCookie(r)
	if err != nil {
		jsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Path[len("/api/notes/"):]
	noteID, err := strconv.Atoi(idStr)
	if err != nil {
		jsonError(w, "invalid note id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// GET /notes/{id} → Any logged-in user can view any note
		var n models.Note
		q := `SELECT n.id, n.owner_id, u.username, n.title, n.content, n.shared, n.favorite, n.updated_at 
		      FROM notes n JOIN users u ON u.id = n.owner_id WHERE n.id=$1`
		err := h.DB.QueryRow(q, noteID).Scan(&n.ID, &n.OwnerID, &n.OwnerUsername, &n.Title, &n.Content, &n.Shared, &n.Favorite, &n.Updated)
		if err != nil {
			jsonError(w, "note not found", http.StatusNotFound)
			return
		}
		jsonResponse(w, n, http.StatusOK)

	case http.MethodPut:
		// PUT /notes/{id} → Any logged-in user can update any note
		var req models.Note
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "invalid json", http.StatusBadRequest)
			return
		}

		_, err = h.DB.Exec(`UPDATE notes 
			SET title=$1, content=$2, shared=$3, favorite=$4, updated_at=now()
			WHERE id=$5`, req.Title, req.Content, req.Shared, req.Favorite, noteID)
		if err != nil {
			jsonError(w, "update failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"message": "updated"}, http.StatusOK)

	case http.MethodDelete:
		// DELETE /notes/{id} → Only owner can delete
		var ownerID int
		err = h.DB.QueryRow(`SELECT owner_id FROM notes WHERE id=$1`, noteID).Scan(&ownerID)
		if err != nil {
			jsonError(w, "note not found", http.StatusNotFound)
			return
		}
		if ownerID != userID {
			jsonError(w, "forbidden: only owner can delete", http.StatusForbidden)
			return
		}
		_, err = h.DB.Exec(`DELETE FROM notes WHERE id=$1`, noteID)
		if err != nil {
			jsonError(w, "delete failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"message": "deleted"}, http.StatusOK)

	default:
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
