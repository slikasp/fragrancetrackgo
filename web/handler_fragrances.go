package web

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/slikasp/fragrancetrackgo/internal/database/remoteDatabase"
)

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
