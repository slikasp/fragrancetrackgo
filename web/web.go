package web

import (
	"context"
	"errors"
	"html"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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
	Brand        string
	Name         string
	Rating       string
	Comment      string
	RatingValue  string
	CommentValue string
}

// shared payload passed into templates
type pageData struct {
	Title      string
	AuthMode   string
	IsAuthed   bool
	UserName   string
	Query      string
	Fragrances []fragranceItem
	MyRatings  []ratingItem
}

// create webapp object
func New(users *localDatabase.Queries, fragrances *remoteDatabase.Queries) (*webApp, error) {
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

// starts the HTTP server and map all handlers
func (a *webApp) Serve(addr string) error {
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

	// User-specific ratings endpoints (authenticated).
	mux.Handle("/my-ratings", a.requireAuth(http.HandlerFunc(a.handleRatingsPage)))
	mux.Handle("/my-ratings/list", a.requireAuth(http.HandlerFunc(a.handleRatingsList)))
	mux.Handle("/my-ratings/add", a.requireAuth(http.HandlerFunc(a.handleRatingsAdd)))
	mux.Handle("/my-ratings/update", a.requireAuth(http.HandlerFunc(a.handleRatingsUpdate)))
	mux.Handle("/my-ratings/remove", a.requireAuth(http.HandlerFunc(a.handleRatingsRemove)))

	log.Printf("web server listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
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

// middleware for protected routes
// It resolves session -> user and injects that user into request context
func (a *webApp) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := a.currentUserFromRequest(r)
		if err != nil {
			// For HTMX requests, use HX-Redirect instead of regular redirect so
			// the browser navigates from within the partial-response flow.
			if isHTMXRequest(r) {
				w.Header().Set("HX-Redirect", "/auth/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
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
	// log.Println(message)
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
