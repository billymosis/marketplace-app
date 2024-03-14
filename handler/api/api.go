package api

import (
	"net/http"

	"github.com/billymosis/marketplace-app/handler/api/product"
	"github.com/billymosis/marketplace-app/handler/api/user"
	appMiddleware "github.com/billymosis/marketplace-app/middleware"
	"github.com/billymosis/marketplace-app/model"
	ps "github.com/billymosis/marketplace-app/store/product"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	Users    model.UserStore
	Products *ps.ProductStore
}

func New(users model.UserStore, products *ps.ProductStore) Server {
	return Server{
		Users:    users,
		Products: products,
	}
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/v1/user", func(r chi.Router) {
		r.Post("/login", user.HandleAuthentication(s.Users))
		r.Post("/register", user.HandleRegistration(s.Users))
	})
	r.Route("/v1/product", func(r chi.Router) {
		r.Use(appMiddleware.ValidateJWT)
		r.Get("/", product.HandleGetProducts(s.Products))
		r.Post("/", product.HandleCreateProduct(s.Products))
		r.Patch("/{id}", product.HandleUpdateProduct(s.Products))
		r.Delete("/{id}", product.HandleDeleteProduct(s.Products))
	})
	return r
}
