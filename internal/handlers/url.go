package handlers 

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/finlleyl/shorty_reborn/internal/service"
)

type Handler struct {
	URLService service.URLService
}

func NewHandler(urlService service.URLService) *Handler {
	return &Handler{URLService: urlService}
}

func (h *Handler) URLRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", h.Create)
	r.Get("/{alias}", h.Resolve)
	r.Delete("/{alias}", h.Delete)

	return r
}

type createURLRequest struct {
	URL string `json:"url"`
	Alias string `json:"alias"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1 << 20)
	defer r.Body.Close()

	var req createURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u, err := h.URLService.Create(r.Context(), req.URL, req.Alias)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAliasExists):
			writeJSONError(w, http.StatusConflict, "alias already exists")
		default:
			writeJSONError(w, http.StatusInternalServerError, "failed to create url")
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/api/urls/%s", u.Alias))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)	
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
  }
  
