package product

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler is the concrete HTTP adapter for the catalogue. It stays thin: it
// decodes the request, calls the Service and renders templates or JSON.
type Handler struct {
	service Service
	tmpl    *template.Template
}

// NewHandler builds the catalogue handler from its service and template set.
func NewHandler(s Service, tmpl *template.Template) *Handler {
	return &Handler{service: s, tmpl: tmpl}
}

// Routes mounts the RESTful product endpoints.
func (h *Handler) Routes(r chi.Router) {
	r.Get("/products", h.list)             // browse / filter (HTML grid)
	r.Post("/products", h.create)          // add a product (JSON)
	r.Get("/products/{uuid}", h.get)       // single product (JSON)
	r.Patch("/products/{uuid}", h.update)  // edit a product (JSON)
	r.Delete("/products/{uuid}", h.delete) // remove a product
}

// View decorates a Product with presentation-only data for templates.
type View struct {
	Product
	SizesJSON template.JS
}

// ToView builds a template view for a single product.
func ToView(p Product) View {
	b, _ := json.Marshal(p.Sizes)
	return View{Product: p, SizesJSON: template.JS(b)}
}

// ToViews builds template views for a slice of products.
func ToViews(ps []Product) []View {
	views := make([]View, len(ps))
	for i, p := range ps {
		views[i] = ToView(p)
	}
	return views
}

// list renders the product grid, optionally filtered by ?family=.
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.List(r.URL.Query().Get("family"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, "product-grid", ToViews(products))
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	p, err := h.service.Get(chi.URLParam(r, "uuid"))
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	created, err := h.service.Create(p)
	if errors.Is(err, ErrInvalid) {
		http.Error(w, "invalid product", http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	updated, err := h.service.Update(chi.URLParam(r, "uuid"), p)
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	err := h.service.Delete(chi.URLParam(r, "uuid"))
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
