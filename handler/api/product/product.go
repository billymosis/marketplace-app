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
	as "github.com/billymosis/marketplace-app/store/account"
	ps "github.com/billymosis/marketplace-app/store/product"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
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
		if id == "" {
			render.NotFound(w, errors.New("not found"))
		}

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

func GetProductDetail(ps *ps.ProductStore, as *as.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		product, err := ps.GetProductById(r.Context(), uint(id))
		if err != nil {
			render.InternalError(w, err)
			return
		}
		accounts, err := as.GetAccountByUser(r.Context(), product.UserId)
		total, err := ps.GetTotalSold(r.Context(), uint(id))
		logrus.Printf("%+v\n", accounts)
		logrus.Printf("%+v\n", total)
		if err != nil {
			render.InternalError(w, err)
			return
		}
		var ba []BankAccount
		for _, element := range accounts {
			ba = append(ba, BankAccount{
				BankAccountID:     strconv.FormatUint(uint64(element.Id), 10),
				BankName:          element.Name,
				BankAccountName:   element.AccountName,
				BankAccountNumber: element.AccountNumber,
			})

		}
		logrus.Printf("%+v\n", product)
		logrus.Printf("%+v\n", ba)
		response := GetProductDetailResponse{
			Message: "ok",
			Data: struct {
				Product ProductResponse `json:"product"`
				Seller  Seller          `json:"seller"`
			}{

				Product: ProductResponse{
					ProductId:     strconv.FormatUint(uint64(product.Id), 10),
					Name:          product.Name,
					Price:         product.Price,
					ImageUrl:      product.ImageUrl,
					Stock:         product.Stock,
					Condition:     product.Condition,
					Tags:          product.Tags,
					IsPurchasable: product.IsPurchasable,
					PurchaseCount: product.PurchaseCount,
				},
				Seller: Seller{
					Name:             "John",
					ProductSoldTotal: total,
					BankAccounts:     ba,
				},
			},
		}

		render.JSON(w, response, http.StatusOK)
	}
}

func UpdateStock(ps *ps.ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			render.NotFound(w, errors.New("not found"))
		}

		productId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		var req updateProductStockRequest
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
			Id:    uint(productId),
			Stock: uint(req.Stock),
		}
		result, err := ps.UpdateProductStock(r.Context(), &product)

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

func Buy(ps *ps.ProductStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			render.NotFound(w, errors.New("not found"))
		}

		productId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		var req buyProductRequest
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

		accountId, err := strconv.ParseUint(req.BankAccountId, 10, 64)
		if err != nil {
			render.BadRequest(w, err)
			return
		}
		payment := model.Payment{
			Id:                   uint(productId),
			AccountId:            uint(accountId),
			ProductId:            uint(productId),
			PaymentProofImageUrl: req.PaymentProofImageUrl,
			Quantity:             uint(req.Quantity),
		}
		err = ps.Payment(r.Context(), &payment)
		if err != nil {
			render.InternalError(w, err)
			return
		}

		w.WriteHeader(200)
	}
}
