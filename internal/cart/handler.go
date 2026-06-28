package cart

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const cookieName = "cart_id"

// Handler is the concrete HTTP adapter for the cart.
type Handler struct {
	service Service
	tmpl    *template.Template
}

// NewHandler builds the cart handler from its service and template set.
func NewHandler(s Service, tmpl *template.Template) *Handler {
	return &Handler{service: s, tmpl: tmpl}
}

// Routes mounts the RESTful cart endpoints.
func (h *Handler) Routes(r chi.Router) {
	r.Get("/cart", h.view)                                   // current cart (HTML)
	r.Post("/cart/items", h.addItem)                         // add a line
	r.Delete("/cart/items/{productUUID}/{ml}", h.removeItem) // drop a line
}

// SessionID returns the cart uuid stored in the visitor's cookie, creating and
// setting an httpOnly cookie when one is not present yet.
func SessionID(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie(cookieName); err == nil && c.Value != "" {
		return c.Value
	}
	id := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return id
}

func (h *Handler) view(w http.ResponseWriter, r *http.Request) {
	c, err := h.service.Get(SessionID(w, r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, "cart-body", c)
}

func (h *Handler) addItem(w http.ResponseWriter, r *http.Request) {
	productUUID := r.FormValue("product_id")
	ml, err := strconv.Atoi(r.FormValue("ml"))
	if productUUID == "" || err != nil {
		http.Error(w, "invalid item", http.StatusBadRequest)
		return
	}
	c, err := h.service.AddItem(SessionID(w, r), productUUID, ml, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.render(w, "cart-body", c)
}

func (h *Handler) removeItem(w http.ResponseWriter, r *http.Request) {
	productUUID := chi.URLParam(r, "productUUID")
	ml, err := strconv.Atoi(chi.URLParam(r, "ml"))
	if productUUID == "" || err != nil {
		http.Error(w, "invalid item", http.StatusBadRequest)
		return
	}
	c, err := h.service.RemoveItem(SessionID(w, r), productUUID, ml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, "cart-body", c)
}

func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
