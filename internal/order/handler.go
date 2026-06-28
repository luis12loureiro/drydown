package order

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luis12loureiro/drydown/internal/cart"
)

// Handler is the concrete HTTP adapter for orders. The checkout endpoint reads
// the visitor's cart, so it depends on the cart service.
type Handler struct {
	service Service
	carts   cart.Service
	tmpl    *template.Template
}

// NewHandler builds the order handler from its service, the cart service and
// the template set.
func NewHandler(s Service, carts cart.Service, tmpl *template.Template) *Handler {
	return &Handler{service: s, carts: carts, tmpl: tmpl}
}

// Routes mounts the RESTful order endpoints.
func (h *Handler) Routes(r chi.Router) {
	r.Get("/orders", h.list)             // all orders (JSON)
	r.Post("/orders", h.checkout)        // place order from the cart (HTML)
	r.Get("/orders/{uuid}", h.get)       // single order (JSON)
	r.Patch("/orders/{uuid}", h.update)  // change status (JSON)
	r.Delete("/orders/{uuid}", h.cancel) // cancel an order
}

// checkout turns the visitor's cart into a pending order and empties the cart.
func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	c, err := h.carts.Get(cart.SessionID(w, r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	o, err := h.service.Place(itemsFromCart(c))
	if errors.Is(err, ErrEmpty) {
		http.Error(w, "cart is empty", http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.carts.Clear(c.UUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, "order-confirmation", o)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	o, err := h.service.Get(chi.URLParam(r, "uuid"))
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	o, err := h.service.SetStatus(chi.URLParam(r, "uuid"), body.Status)
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, o)
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	err := h.service.Cancel(chi.URLParam(r, "uuid"))
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

// itemsFromCart maps cart lines into order lines, snapshotting price and name.
func itemsFromCart(c cart.Cart) []Item {
	items := make([]Item, len(c.Items))
	for i, it := range c.Items {
		items[i] = Item{
			ProductUUID: it.ProductUUID,
			Name:        it.Name,
			ML:          it.ML,
			Price:       it.Price,
			Qty:         it.Qty,
		}
	}
	return items
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
