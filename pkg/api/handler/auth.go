package handler

import (
	"encoding/json"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/auth"
	"net/http"
)

// SignUp godoc
// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} utility.ErrorResponse
// @Failure 500 {object} utility.ErrorResponse
// @Failure default {object} utility.ErrorResponse
// @Router /auth/sign-up [post]
func SignUp(service auth.AuthorizationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input todo.User

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		id, err := service.CreateUser(input)
		if err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"id": id,
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

type signInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignIn godoc
// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} utility.ErrorResponse
// @Failure 500 {object} utility.ErrorResponse
// @Failure default {object} utility.ErrorResponse
// @Router /auth/sign-in [post]
func SignIn(service auth.AuthorizationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input signInInput

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		token, err := service.GenerateToken(input.Username, input.Password)
		if err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"token": token,
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
