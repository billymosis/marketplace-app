package api

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/billymosis/marketplace-app/handler/api/account"
	"github.com/billymosis/marketplace-app/handler/api/product"
	"github.com/billymosis/marketplace-app/handler/api/user"
	appMiddleware "github.com/billymosis/marketplace-app/middleware"
	"github.com/billymosis/marketplace-app/model"
	"github.com/billymosis/marketplace-app/service/image"
	as "github.com/billymosis/marketplace-app/store/account"
	ps "github.com/billymosis/marketplace-app/store/product"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	Users    model.UserStore
	Products *ps.ProductStore
	Accounts *as.AccountStore
	S3Client *s3.Client
}

func New(users model.UserStore, products *ps.ProductStore, accounts *as.AccountStore, s3client *s3.Client) Server {
	return Server{
		Users:    users,
		Products: products,
		Accounts: accounts,
		S3Client: s3client,
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
			r.Post("/{id}/stock", product.UpdateStock(s.Products))
			r.Post("/{id}/buy", product.Buy(s.Products))
		})
	})
	r.Route("/v1/account", func(r chi.Router) {
		r.Use(appMiddleware.ValidateJWT)
		r.Post("/", account.Create(s.Accounts))
		r.Patch("/{id}", account.Update(s.Accounts))
		r.Delete("/{id}", account.Delete(s.Accounts))
		r.Get("/", account.Get(s.Accounts))
	})
	r.Route("/v1/image", func(r chi.Router) {
		r.Use(appMiddleware.ValidateJWT)
		r.Post("/", image.Upload(s.S3Client))
	})
	return r
}
