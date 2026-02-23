package web

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/slikasp/fragrancetrackgo/handlers"
)

// logged-in users go to /fragrances, everyone else goes to /auth/login.
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
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

// keeps /auth as a stable alias and forwards unauthenticated users to login.
func (a *webApp) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if _, err := a.currentUserFromRequest(r); err == nil {
		http.Redirect(w, r, "/fragrances", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

// validate form input and register user
func (a *webApp) handleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if _, err := a.currentUserFromRequest(r); err == nil {
			http.Redirect(w, r, "/fragrances", http.StatusSeeOther)
			return
		}
		a.render(w, "auth_page", pageData{
			Title:    "Register",
			AuthMode: "register",
		})
		return
	case http.MethodPost:
		// Continue with registration submit handling below.
	default:
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
	switch r.Method {
	case http.MethodGet:
		if _, err := a.currentUserFromRequest(r); err == nil {
			http.Redirect(w, r, "/fragrances", http.StatusSeeOther)
			return
		}
		a.render(w, "auth_page", pageData{
			Title:    "Login",
			AuthMode: "login",
		})
		return
	case http.MethodPost:
		// Continue with login submit handling below.
	default:
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
		w.Header().Set("HX-Redirect", "/auth/login")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
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
