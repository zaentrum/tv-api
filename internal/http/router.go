package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/zaentrum/tv-api/internal/auth"
	"github.com/zaentrum/tv-api/internal/config"
)

func NewRouter(cfg config.Config) (http.Handler, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/api/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "product": "tv"})
	})

	r.Get("/api/openapi.yaml", serveOpenAPI)

	verifier := auth.NewVerifier(cfg.OIDCIssuer, cfg.OIDCAudience, cfg.OIDCEnabled)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(verifier.Middleware)
		r.Get("/me", whoAmI)
		r.Get("/items", listItems)
	})

	return r, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func whoAmI(w http.ResponseWriter, r *http.Request) {
	sub, err := auth.SubjectFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"sub": sub})
}

// listItems is a stub that returns sample data. Replace with a real catalog
// query once that backend is online.
func listItems(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"product": "tv",
		"items":   sampleItems(),
	})
}
