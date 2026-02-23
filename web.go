package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/slikasp/fragrancetrackgo/handlers"
	"github.com/slikasp/fragrancetrackgo/internal/database/localDatabase"
	"github.com/slikasp/fragrancetrackgo/internal/database/remoteDatabase"
)

const sessionCookieName = "ftg_session"

type contextKey string

const currentUserContextKey contextKey = "current_user"

type webApp struct {
	users      *localDatabase.Queries
	fragrances *remoteDatabase.Queries
	templates  *template.Template
}

// view model for template rendering
type fragranceItem struct {
	Brand string
	Name  string
	URL   string
}

// UI-safe shape of a user's saved fragrance rating
// Nullable DB fields are converted to display strings before rendering
type ratingItem struct {
	Brand   string
	Name    string
	Rating  string
	Comment string
}

// shared payload passed into templates
type pageData struct {
	Title        string
	IsAuthed     bool
	UserName     string
	Query        string
	Fragrances   []fragranceItem
	MyFragrances []ratingItem
}

// create webapp object
func newWebApp(users *localDatabase.Queries, fragrances *remoteDatabase.Queries) (*webApp, error) {
	tmpl, err := loadTemplates()
	if err != nil {
		return nil, err
	}

	a := webApp{
		users:      users,
		fragrances: fragrances,
		templates:  tmpl,
	}

	return &a, nil
}

// load and parse HTML templates
func loadTemplates() (*template.Template, error) {
	var files []string
	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("no template files found in templates/")
	}
	return template.ParseFiles(files...)
}

// starts the HTTP server and map all handlers
func (a *webApp) serve(addr string) error {
	mux := http.NewServeMux()

	// Authentication/session endpoints.
	mux.HandleFunc("/", a.handleHome)
	mux.HandleFunc("/auth", a.handleAuthPage)
	mux.HandleFunc("/auth/login", a.handleLogin)
	mux.HandleFunc("/auth/register", a.handleRegister)
	mux.HandleFunc("/auth/logout", a.handleLogout)

	// Fragrance browsing endpoints (authenticated).
	mux.Handle("/fragrances", a.requireAuth(http.HandlerFunc(a.handleFragrancesPage)))
	mux.Handle("/fragrances/search", a.requireAuth(http.HandlerFunc(a.handleFragrancesSearch)))

	// User-specific collection/rating endpoints (authenticated).
	mux.Handle("/my-fragrances", a.requireAuth(http.HandlerFunc(a.handleMyFragrancesPage)))
	mux.Handle("/my-fragrances/list", a.requireAuth(http.HandlerFunc(a.handleMyFragrancesList)))
	mux.Handle("/my-fragrances/add", a.requireAuth(http.HandlerFunc(a.handleMyFragrancesAdd)))

	log.Printf("web server listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

// PAGE HANDLERS

// logged-in users go to /fragrances, everyone else goes to /auth.
func (a *webApp) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := a.currentUserFromRequest(r)
	if err == nil {
		http.Redirect(w, r, "/fragrances", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

// If the user is already authenticated, redirects to /fragrances.
func (a *webApp) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if _, err := a.currentUserFromRequest(r); err == nil {
		http.Redirect(w, r, "/fragrances", http.StatusSeeOther)
		return
	}
	a.render(w, "auth_page", pageData{Title: "Login / Register"})
}

// validate form input and register user
func (a *webApp) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		a.respondWithError(w, r, "invalid form payload", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	password := r.FormValue("password")

	if name == "" || password == "" {
		a.respondWithError(w, r, "name and password are required", http.StatusBadRequest)
		return
	}

	user, err := handlers.UserRegister(r.Context(), a.users, name, password)
	if err != nil {
		// Map domain errors to stable HTTP responses.
		switch {
		case errors.Is(err, handlers.ErrInvalidInput):
			a.respondWithError(w, r, "name and password are required", http.StatusBadRequest)
		case errors.Is(err, handlers.ErrUserAlreadyExists):
			a.respondWithError(w, r, "user already exists", http.StatusConflict)
		default:
			log.Printf("register failed: %v", err)
			a.respondWithError(w, r, "failed to register user", http.StatusInternalServerError)
		}
		return
	}

	a.setSessionCookie(w, user.ID)
	a.respondWithRedirectOrSuccess(w, r, "/fragrances", "Registration successful.")
}

// validate login credentials
func (a *webApp) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		a.respondWithError(w, r, "invalid form payload", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	password := r.FormValue("password")
	if name == "" || password == "" {
		a.respondWithError(w, r, "name and password are required", http.StatusBadRequest)
		return
	}

	user, err := handlers.UserLogin(r.Context(), a.users, name, password)
	if err != nil {
		// Domain-layer sentinel errors keep handler logic clear and predictable.
		switch {
		case errors.Is(err, handlers.ErrInvalidInput):
			a.respondWithError(w, r, "name and password are required", http.StatusBadRequest)
		case errors.Is(err, handlers.ErrInvalidCredential):
			a.respondWithError(w, r, "invalid credentials", http.StatusUnauthorized)
		default:
			log.Printf("login failed: %v", err)
			a.respondWithError(w, r, "failed to login", http.StatusInternalServerError)
		}
		return
	}

	a.setSessionCookie(w, user.ID)
	a.respondWithRedirectOrSuccess(w, r, "/fragrances", "Login successful.")
}

// clears the cookie and returns either:
// - HX-Redirect for HTMX requests (partial page context), or
// - normal browser redirect for full-page form submits.
func (a *webApp) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	a.clearSessionCookie(w)
	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", "/auth")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

// handleFragrancesPage renders the full "all fragrances" page.
func (a *webApp) handleFragrancesPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	frags, err := a.searchFragrances(r.Context(), q, 50, 0)
	if err != nil {
		log.Printf("search fragrances failed: %v", err)
		http.Error(w, "failed to load fragrances", http.StatusInternalServerError)
		return
	}

	a.render(w, "fragrances_page", pageData{
		Title:      "All Fragrances",
		IsAuthed:   true,
		UserName:   user.Name,
		Query:      q,
		Fragrances: frags,
	})
}

// returns only table/list rows for HTMX updates
// shares search logic with handleFragrancesPage but renders a partial template
func (a *webApp) handleFragrancesSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	frags, err := a.searchFragrances(r.Context(), q, 50, 0)
	if err != nil {
		log.Printf("search fragrances failed: %v", err)
		http.Error(w, "failed to search fragrances", http.StatusInternalServerError)
		return
	}

	a.render(w, "fragrance_rows", pageData{Fragrances: frags})
}

// calls the remote fragrance catalog DB
func (a *webApp) searchFragrances(ctx context.Context, q string, limit, offset int32) ([]fragranceItem, error) {
	rows, err := a.fragrances.SearchFragrances(ctx, remoteDatabase.SearchFragrancesParams{
		Btrim:  q,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	frags := make([]fragranceItem, 0, len(rows))
	for _, row := range rows {
		item := fragranceItem{}
		if row.Brand.Valid {
			item.Brand = row.Brand.String
		}
		if row.Name.Valid {
			item.Name = row.Name.String
		}
		if row.Url.Valid {
			item.URL = row.Url.String
		}
		frags = append(frags, item)
	}
	return frags, nil
}

// renders the full "my fragrances" page for the current user
func (a *webApp) handleMyFragrancesPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ratings, err := a.listMyFragrances(r.Context(), user.ID)
	if err != nil {
		log.Printf("list my fragrances failed: %v", err)
		http.Error(w, "failed to load your fragrances", http.StatusInternalServerError)
		return
	}

	a.render(w, "my_fragrances_page", pageData{
		Title:        "My Fragrances",
		IsAuthed:     true,
		UserName:     user.Name,
		MyFragrances: ratings,
	})
}

// returns only row markup used by HTMX to refresh the list
func (a *webApp) handleMyFragrancesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := userFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ratings, err := a.listMyFragrances(r.Context(), user.ID)
	if err != nil {
		log.Printf("list my fragrances failed: %v", err)
		http.Error(w, "failed to load your fragrances", http.StatusInternalServerError)
		return
	}

	a.render(w, "my_fragrance_rows", pageData{MyFragrances: ratings})
}

// validates add-form data and adds a fragrance from the public list to my list
func (a *webApp) handleMyFragrancesAdd(w http.ResponseWriter, r *http.Request) {
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
		_, _ = w.Write([]byte(`<div class="ok">Added to your fragrances.</div>`))
		return
	}

	http.Redirect(w, r, "/my-fragrances", http.StatusSeeOther)
}

// reads saved ratings and converts nullable DB fields to display values
// The UI uses "-" for missing rating/comment
func (a *webApp) listMyFragrances(ctx context.Context, userID uuid.UUID) ([]ratingItem, error) {
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
		} else {
			item.Rating = "-"
		}
		if row.Comment.Valid {
			item.Comment = row.Comment.String
		} else {
			item.Comment = "-"
		}
		ratings = append(ratings, item)
	}
	return ratings, nil
}

// middleware for protected routes
// It resolves session -> user and injects that user into request context
func (a *webApp) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := a.currentUserFromRequest(r)
		if err != nil {
			// For HTMX requests, use HX-Redirect instead of regular redirect so
			// the browser navigates from within the partial-response flow.
			if isHTMXRequest(r) {
				w.Header().Set("HX-Redirect", "/auth")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), currentUserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// performs the full session lookup:
// read session cookie -> parse UUID -> query users table.
func (a *webApp) currentUserFromRequest(r *http.Request) (localDatabase.User, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return localDatabase.User{}, err
	}
	userID, err := uuid.Parse(cookie.Value)
	if err != nil {
		return localDatabase.User{}, err
	}
	return a.users.GetUserByID(r.Context(), userID)
}

// reads the user object injected by requireAuth middleware
func userFromContext(ctx context.Context) (localDatabase.User, bool) {
	user, ok := ctx.Value(currentUserContextKey).(localDatabase.User)
	return user, ok
}

// stores logged-in user ID in an HttpOnly cookie
// SameSite=Lax helps reduce CSRF risk for cross-site requests
func (a *webApp) setSessionCookie(w http.ResponseWriter, userID uuid.UUID) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    userID.String(),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// invalidates the session cookie in the browser
func (a *webApp) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// executes a named HTML template with standard content-type handling
func (a *webApp) render(w http.ResponseWriter, templateName string, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := a.templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		log.Printf("template render failed (%s): %v", templateName, err)
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}
}

// normalizes post-success behavior for HTMX/non-HTMX clients
func (a *webApp) respondWithRedirectOrSuccess(w http.ResponseWriter, r *http.Request, redirectPath, message string) {
	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", redirectPath)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
	if message != "" {
		log.Println(message)
	}
}

// normalizes error responses for HTMX/non-HTMX clients
// HTMX gets a fragment that can be swapped into the page; full requests get http.Error
func (a *webApp) respondWithError(w http.ResponseWriter, r *http.Request, message string, status int) {
	if isHTMXRequest(r) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`<div class="err">` + html.EscapeString(message) + `</div>`))
		return
	}
	http.Error(w, message, status)
}

// identifies requests initiated by HTMX
func isHTMXRequest(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("HX-Request"), "true")
}
