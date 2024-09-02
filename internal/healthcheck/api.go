package healthcheck

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterHandlers(r *chi.Mux) {
	r.Get("/healthcheck", healthCheckHandler)
}

// healthcheck responds to a healthcheck request.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
