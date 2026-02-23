package web

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/slikasp/fragrancetrackgo/handlers"
)

// renders the full "my fragrances" page for the current user
func (a *webApp) handleRatingsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ratings, err := a.listRatings(r.Context(), user.ID)
	if err != nil {
		log.Printf("list my ratings failed: %v", err)
		http.Error(w, "failed to load your ratings", http.StatusInternalServerError)
		return
	}

	a.render(w, "my_ratings_page", pageData{
		Title:     "My ratings",
		IsAuthed:  true,
		UserName:  user.Name,
		MyRatings: ratings,
	})
}

// returns only row markup used by HTMX to refresh the list
func (a *webApp) handleRatingsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ratings, err := a.listRatings(r.Context(), user.ID)
	if err != nil {
		log.Printf("list my ratings failed: %v", err)
		http.Error(w, "failed to load your ratings", http.StatusInternalServerError)
		return
	}

	a.render(w, "my_rating_rows", pageData{MyRatings: ratings})
}

// validates add-form data and adds a fragrance from the public list to my list
func (a *webApp) handleRatingsAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		a.respondWithError(w, r, "invalid form payload", http.StatusBadRequest)
		return
	}

	brand := strings.ToLower(strings.TrimSpace(r.FormValue("brand")))
	name := strings.ToLower(strings.TrimSpace(r.FormValue("name")))
	commentRaw := strings.TrimSpace(r.FormValue("comment"))
	ratingRaw := strings.TrimSpace(r.FormValue("rating"))

	if brand == "" || name == "" {
		a.respondWithError(w, r, "brand and fragrance name are required", http.StatusBadRequest)
		return
	}

	var rating sql.NullInt32
	if ratingRaw != "" {
		// Keep rating constraints close to the HTTP boundary for clear user feedback.
		v, err := strconv.Atoi(ratingRaw)
		if err != nil || v < 0 || v > 10 {
			a.respondWithError(w, r, "rating must be an integer between 0 and 10", http.StatusBadRequest)
			return
		}
		rating = sql.NullInt32{Int32: int32(v), Valid: true}
	}

	var comment sql.NullString
	if commentRaw != "" {
		comment = sql.NullString{String: commentRaw, Valid: true}
	}

	_, err := handlers.RatingAdd(r.Context(), a.users, user.ID, handlers.RatingInput{
		Brand:   brand,
		Name:    name,
		Rating:  rating,
		Comment: comment,
	})
	if err != nil {
		switch {
		case errors.Is(err, handlers.ErrInvalidInput):
			a.respondWithError(w, r, "brand and fragrance name are required", http.StatusBadRequest)
		case errors.Is(err, handlers.ErrRatingExists):
			a.respondWithError(w, r, "fragrance already exists in your collection", http.StatusConflict)
		default:
			log.Printf("add rating failed: %v", err)
			a.respondWithError(w, r, "failed to save fragrance", http.StatusInternalServerError)
		}
		return
	}

	if isHTMXRequest(r) {
		// Trigger an event listened to by the table body so it can refetch rows.
		w.Header().Set("HX-Trigger", "fragranceAdded")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`<div class="ok">Added to your ratings.</div>`))
		return
	}

	http.Redirect(w, r, "/my-ratings", http.StatusSeeOther)
}

// validates edit-form data and updates score/comment for an existing fragrance
func (a *webApp) handleRatingsUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		a.respondWithError(w, r, "invalid form payload", http.StatusBadRequest)
		return
	}

	brand := strings.ToLower(strings.TrimSpace(r.FormValue("brand")))
	name := strings.ToLower(strings.TrimSpace(r.FormValue("name")))
	commentRaw := strings.TrimSpace(r.FormValue("comment"))
	ratingRaw := strings.TrimSpace(r.FormValue("rating"))

	if brand == "" || name == "" {
		a.respondWithError(w, r, "brand and fragrance name are required", http.StatusBadRequest)
		return
	}

	var rating sql.NullInt32
	if ratingRaw != "" {
		v, err := strconv.Atoi(ratingRaw)
		if err != nil || v < 0 || v > 10 {
			a.respondWithError(w, r, "rating must be an integer between 0 and 10", http.StatusBadRequest)
			return
		}
		rating = sql.NullInt32{Int32: int32(v), Valid: true}
	}

	var comment sql.NullString
	if commentRaw != "" {
		comment = sql.NullString{String: commentRaw, Valid: true}
	}

	_, err := handlers.RatingUpdate(r.Context(), a.users, user.ID, brand, name, comment, rating)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			a.respondWithError(w, r, "fragrance not found in your collection", http.StatusNotFound)
		default:
			log.Printf("update rating failed: %v", err)
			a.respondWithError(w, r, "failed to update fragrance", http.StatusInternalServerError)
		}
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Trigger", "fragranceUpdated")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<div class="ok">Fragrance updated.</div>`))
		return
	}

	http.Redirect(w, r, "/my-ratings", http.StatusSeeOther)
}

// validates remove request and deletes one saved fragrance for current user
func (a *webApp) handleRatingsRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		a.respondWithError(w, r, "invalid form payload", http.StatusBadRequest)
		return
	}

	brand := strings.ToLower(strings.TrimSpace(r.FormValue("brand")))
	name := strings.ToLower(strings.TrimSpace(r.FormValue("name")))
	if brand == "" || name == "" {
		a.respondWithError(w, r, "brand and fragrance name are required", http.StatusBadRequest)
		return
	}

	_, err := handlers.RatingRemove(r.Context(), a.users, user.ID, brand, name)
	if err != nil {
		switch {
		case errors.Is(err, handlers.ErrInvalidInput):
			a.respondWithError(w, r, "brand and fragrance name are required", http.StatusBadRequest)
		case errors.Is(err, sql.ErrNoRows):
			a.respondWithError(w, r, "fragrance not found in your collection", http.StatusNotFound)
		default:
			log.Printf("remove rating failed: %v", err)
			a.respondWithError(w, r, "failed to remove fragrance", http.StatusInternalServerError)
		}
		return
	}

	if isHTMXRequest(r) {
		w.Header().Set("HX-Trigger", "fragranceRemoved")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<div class="ok">Fragrance removed.</div>`))
		return
	}

	http.Redirect(w, r, "/my-ratings", http.StatusSeeOther)
}

// reads saved ratings and converts nullable DB fields to display values
// The UI uses "-" for missing rating/comment
func (a *webApp) listRatings(ctx context.Context, userID uuid.UUID) ([]ratingItem, error) {
	rows, err := handlers.RatingList(ctx, a.users, userID)
	if err != nil {
		return nil, err
	}

	ratings := make([]ratingItem, 0, len(rows))
	for _, row := range rows {
		item := ratingItem{
			Brand: row.Brand,
			Name:  row.Name,
		}
		if row.Rating.Valid {
			item.Rating = fmt.Sprintf("%d", row.Rating.Int32)
			item.RatingValue = item.Rating
		} else {
			item.Rating = "-"
			item.RatingValue = ""
		}
		if row.Comment.Valid {
			item.Comment = row.Comment.String
			item.CommentValue = row.Comment.String
		} else {
			item.Comment = "-"
			item.CommentValue = ""
		}
		ratings = append(ratings, item)
	}
	return ratings, nil
}
