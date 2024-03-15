package account

import (
	// "database/sql"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	// "errors"
	"io"
	"net/http"

	// "strconv"

	"github.com/billymosis/marketplace-app/handler/render"
	"github.com/billymosis/marketplace-app/model"
	as "github.com/billymosis/marketplace-app/store/account"
	"github.com/go-chi/chi"
	// "github.com/go-chi/chi"
)

func Create(as *as.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req createAccountRequest
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

		if err := as.Validate.Struct(req); err != nil {
			render.BadRequest(w, err)
			return
		}
		account := model.Account{
			Name:          req.BankName,
			AccountName:   req.BankAccountName,
			AccountNumber: req.BankAccountNumber,
		}
		result, err := as.Create(r.Context(), &account)

		if result == nil {
			render.InternalError(w, err)
			return
		}
		w.WriteHeader(200)
	}
}

func Update(as *as.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		accountId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		var req createAccountRequest
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

		if err := as.Validate.Struct(req); err != nil {
			render.BadRequest(w, err)
			return
		}

		account := model.Account{
			Name:          req.BankName,
			AccountName:   req.BankAccountName,
			AccountNumber: req.BankAccountNumber,
			Id:            uint(accountId),
		}
		result, err := as.Update(r.Context(), &account)

		if result == nil {
			if errors.Is(err, sql.ErrNoRows) {
				render.NotFound(w, errors.New("account not found"))
				return
			}
			render.InternalError(w, err)
			return
		}
		w.WriteHeader(200)
	}
}

func Delete(as *as.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		productId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		err = as.Delete(r.Context(), uint(productId))

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				render.NotFound(w, errors.New("account not found"))
				return
			}
			render.InternalError(w, err)
			return
		}
		w.WriteHeader(200)
	}
}
func Get(as *as.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		accounts, err := as.Get(r.Context())
		if err != nil {
			render.InternalError(w, err)
			return
		}

		var ba []BankAccount
		for _, element := range accounts {
			ba = append(ba, BankAccount{
				BankAccountId:     strconv.FormatUint(uint64(element.Id), 10),
				BankName:          element.Name,
				BankAccountName:   element.AccountName,
				BankAccountNumber: element.AccountNumber,
			})
		}

		response := GetAccountResponse{
			Message: "success",
			Data:    ba,
		}

		render.JSON(w, response, http.StatusOK)
	}
}
