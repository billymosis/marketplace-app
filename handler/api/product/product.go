package product

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/billymosis/marketplace-app/handler/render"
	"github.com/billymosis/marketplace-app/model"
	ps "github.com/billymosis/marketplace-app/store/product"
	"github.com/go-chi/chi"
)

func HandleCreateProduct(ps *ps.ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req createProductRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			render.BadRequest(w, err)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &req); err != nil {
			render.BadRequest(w, err)
			return
		}

		if err := ps.Validate.Struct(req); err != nil {
			render.BadRequest(w, err)
			return
		}
		product := model.Product{
			Name:          req.Name,
			Tags:          req.Tags,
			Price:         uint(req.Price),
			Stock:         uint(req.Stock),
			ImageUrl:      req.ImageURL,
			Condition:     req.Condition,
			IsPurchasable: req.IsPurchasable,
		}
		result, err := ps.CreateProduct(r.Context(), &product)

		if result == nil {
			render.InternalError(w, err)
			return
		}
		w.WriteHeader(200)
	}
}

func HandleUpdateProduct(ps *ps.ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		productId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		var req createProductRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			render.BadRequest(w, err)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &req); err != nil {
			render.BadRequest(w, err)
			return
		}

		if err := ps.Validate.Struct(req); err != nil {
			render.BadRequest(w, err)
			return
		}
		product := model.Product{
			Id:            uint(productId),
			Name:          req.Name,
			Tags:          req.Tags,
			Price:         uint(req.Price),
			Stock:         uint(req.Stock),
			ImageUrl:      req.ImageURL,
			Condition:     req.Condition,
			IsPurchasable: req.IsPurchasable,
		}
		result, err := ps.UpdateProduct(r.Context(), &product)

		if result == nil {
			if errors.Is(err, sql.ErrNoRows) {
				render.NotFound(w, errors.New("product not found"))
				return
			}
			render.InternalError(w, err)
			return
		}
		w.WriteHeader(200)
	}
}

func HandleDeleteProduct(ps *ps.ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		productId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		err = ps.DeleteProduct(r.Context(), uint(productId))

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				render.NotFound(w, errors.New("product not found"))
				return
			}
			render.InternalError(w, err)
			return
		}
		w.WriteHeader(200)
	}
}

func HandleGetProducts(ps *ps.ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		products, meta, err := ps.GetProducts(r.Context(), r.URL.Query())
		if err != nil {
			render.InternalError(w, err)
			return
		}
		var pr []ProductResponse

		for _, element := range products {
			pr = append(pr, ProductResponse{
				ProductId:     strconv.FormatUint(uint64(element.Id), 10),
				Name:          element.Name,
				Price:         element.Price,
				ImageUrl:      element.ImageUrl,
				Stock:         element.Stock,
				Condition:     element.Condition,
				Tags:          element.Tags,
				IsPurchasable: element.IsPurchasable,
				PurchaseCount: element.PurchaseCount,
			})
		}

		response := GetProductsResponse{
			Message: "ok",
			Data:    pr,
			Meta: struct {
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
				Total  int `json:"total"`
			}{
				Limit:  meta.Limit,
				Offset: meta.Offset,
				Total:  meta.Total,
			},
		}

		render.JSON(w, response, http.StatusOK)
	}
}
