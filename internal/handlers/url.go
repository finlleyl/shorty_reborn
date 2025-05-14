package handlers 

import (
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

