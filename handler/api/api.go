package api

import (
	"net/http"

	"github.com/billymosis/marketplace-app/handler/api/account"
	"github.com/billymosis/marketplace-app/handler/api/product"
	"github.com/billymosis/marketplace-app/handler/api/user"
	appMiddleware "github.com/billymosis/marketplace-app/middleware"
	"github.com/billymosis/marketplace-app/model"
	as "github.com/billymosis/marketplace-app/store/account"
	ps "github.com/billymosis/marketplace-app/store/product"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	Users    model.UserStore
	Products *ps.ProductStore
	Accounts *as.AccountStore
}

func New(users model.UserStore, products *ps.ProductStore, accounts *as.AccountStore) Server {
	return Server{
		Users:    users,
		Products: products,
		Accounts: accounts,
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
		r.Get("/", product.HandleGetProducts(s.Products))
		r.Route("/", func(r chi.Router) {
			r.Use(appMiddleware.ValidateJWT)
			r.Post("/", product.HandleCreateProduct(s.Products))
			r.Patch("/{id}", product.HandleUpdateProduct(s.Products))
			r.Delete("/{id}", product.HandleDeleteProduct(s.Products))
		})
	})
	r.Route("/v1/account", func(r chi.Router) {
		r.Use(appMiddleware.ValidateJWT)
		r.Post("/", account.Create(s.Accounts))
		r.Patch("/{id}", account.Update(s.Accounts))
		r.Delete("/{id}", account.Delete(s.Accounts))
		r.Get("/", account.Get(s.Accounts))
	})

	return r
}
