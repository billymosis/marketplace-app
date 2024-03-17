package api

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/billymosis/marketplace-app/handler/api/account"
	"github.com/billymosis/marketplace-app/handler/api/product"
	"github.com/billymosis/marketplace-app/handler/api/user"
	"github.com/billymosis/marketplace-app/middleware"
	"github.com/billymosis/marketplace-app/model"
	"github.com/billymosis/marketplace-app/service/image"
	as "github.com/billymosis/marketplace-app/store/account"
	ps "github.com/billymosis/marketplace-app/store/product"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
func prometheusHandler() http.Handler {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	handler := promhttp.Handler()
	return handler
}

func (s Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/v1", func(r chi.Router) {
		r.Use(AppMiddleware.WrapWithPrometheus)

		r.Route("/user", func(r chi.Router) {
			r.Post("/login", user.HandleAuthentication(s.Users))
			r.Post("/register", user.HandleRegistration(s.Users))
		})

		r.Route("/product", func(r chi.Router) {
			r.Get("/", product.HandleGetProducts(s.Products))
			r.Get("/{id}", product.GetProductDetail(s.Products, s.Accounts))
			r.Route("/", func(r chi.Router) {
				r.Use(AppMiddleware.ValidateJWT)
				r.Post("/", product.HandleCreateProduct(s.Products))
				r.Patch("/{id}", product.HandleUpdateProduct(s.Products))
				r.Delete("/{id}", product.HandleDeleteProduct(s.Products))
				r.Post("/{id}/stock", product.UpdateStock(s.Products))
				r.Post("/{id}/buy", product.Buy(s.Products))
			})
		})
		r.Route("/bank/account", func(r chi.Router) {
			r.Use(AppMiddleware.ValidateJWT)
			r.Post("/", account.Create(s.Accounts))
			r.Get("/", account.Get(s.Accounts))
			r.Patch("/", account.Update(s.Accounts))
			r.Patch("/{id}", account.Update(s.Accounts))
			r.Delete("/{id}", account.Delete(s.Accounts))
		})
	})

	r.Route("/v1/image", func(r chi.Router) {
		r.Use(AppMiddleware.ValidateJWT)
		r.Post("/", image.Upload(s.S3Client))
	})
	return r
}
