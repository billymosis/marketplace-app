package user

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/billymosis/marketplace-app/handler/render"
	"github.com/billymosis/marketplace-app/model"
	"github.com/billymosis/marketplace-app/service/auth"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func HandleAuthentication(us model.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginUserRequest

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

		if err := us.GetValidator().Struct(req); err != nil {
			render.BadRequest(w, err)
			return
		}

		user, err := us.GetByUsername(req.Username)

		if err != nil {
			render.NotFound(w, errors.New("User not found"))
			logrus.Info("api: cannot find user")
			return
		}

		validUser := user.CheckPassword(req.Password)
		if !validUser {
			render.BadRequest(w, errors.New("Invalid username or password"))
			return

		}

		token, err := auth.GenerateToken(user.Id, user.Username)
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		var res loginUserResponse
		res.Message = "User logged successfully"
		res.Data.Name = user.Name
		res.Data.Username = user.Username
		res.Data.AccessToken = token

		render.JSON(w, res, http.StatusOK)
	}
}

func HandleRegistration(us model.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req createUserRequest
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

		if err := us.GetValidator().Struct(req); err != nil {
			render.BadRequest(w, err)
			return
		}

		user := model.User{
			Username: req.Username,
			Name:     req.Name,
			Password: req.Password,
		}
		err = user.HashPassword()
		if err != nil {
			render.BadRequest(w, err)
			return
		}

		result, err := us.CreateUser(r.Context(), &user)
		if err != nil {
			render.ErrorCode(w, errors.New("username already exist"), 409)
			return
		}

		if result == nil {
			render.InternalError(w, err)
			return
		}

		token, err := auth.GenerateToken(result.Id, result.Username)
		if err != nil {
			render.InternalError(w, err)
			return
		}

		var res createUserResponse
		res.Message = "User registered successfully"
		res.Data.Username = user.Username
		res.Data.Name = user.Name
		res.Data.AccessToken = token

		render.JSON(w, res, http.StatusCreated)
	}
}
