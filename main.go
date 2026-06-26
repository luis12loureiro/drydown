// Command drydown is the HTTP entrypoint and composition root: it wires the
// in-memory adapters, domain services and handlers together and serves the app.
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/luis12loureiro/drydown/internal/cart"
	"github.com/luis12loureiro/drydown/internal/order"
	"github.com/luis12loureiro/drydown/internal/product"
	"github.com/luis12loureiro/drydown/internal/user"
)

func main() {
	tmpl := template.Must(template.New("").ParseGlob("templates/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("templates/partials/*.html"))

	// Adapters (in-memory for now; swap for sqlite.go later).
	productRepo := product.NewMemoryRepository()
	cartRepo := cart.NewMemoryRepository()
	orderRepo := order.NewMemoryRepository()
	userRepo := user.NewMemoryRepository()

	// Services (business logic).
	productSvc := product.NewService(productRepo)
	cartSvc := cart.NewService(cartRepo, productSvc)
	orderSvc := order.NewService(orderRepo)
	userSvc := user.NewService(userRepo)

	// Handlers (HTTP adapters).
	productH := product.NewHandler(productSvc, tmpl)
	cartH := cart.NewHandler(cartSvc, tmpl)
	orderH := order.NewHandler(orderSvc, cartSvc, tmpl)
	userH := user.NewHandler(userSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static assets.
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Landing page.
	r.Get("/", home(productSvc, cartSvc, tmpl))

	// Domain routes.
	productH.Routes(r)
	cartH.Routes(r)
	orderH.Routes(r)
	userH.Routes(r)

	addr := ":" + port()
	log.Printf("Drydown listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

// home renders the landing page with the catalogue and the visitor's cart.
func home(products product.Service, carts cart.Service, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ps, err := products.List("")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c, err := carts.Get(cart.SessionID(w, r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Products []product.View
			Cart     cart.Cart
		}{
			Products: product.ToViews(ps),
			Cart:     c,
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
