package handler

import (
	"encoding/json"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/list"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// CreateList godoc
// @Summary Create todo list
// @Security ApiKeyAuth
// @Tags lists
// @Description create todo list
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body todo.TodoList true "list info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} utility.ErrorResponse
// @Failure 500 {object} utility.ErrorResponse
// @Failure default {object} utility.ErrorResponse
// @Router /api/lists [post]
func CreateList(ctx utility.AppContext, service list.TodoListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		var input todo.TodoList
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		id, err := service.Create(ctx, userId, input)
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

type getAllListsResponse struct {
	Data []todo.TodoList `json:"data"`
}

// GetAllLists godoc
// @Summary Get All Lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all lists
// @ID get-all-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllListsResponse
// @Failure 400,404 {object} utility.ErrorResponse
// @Failure 500 {object} utility.ErrorResponse
// @Failure default {object} utility.ErrorResponse
// @Router /api/lists [get]
func GetAllLists(ctx utility.AppContext, service list.TodoListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		lists, err := service.GetAll(ctx, userId)
		if err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := getAllListsResponse{
			Data: lists,
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// GetListById godoc
// @Summary Get List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description get list by id
// @ID get-list-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} todo.ListsItem
// @Failure 400,404 {object} utility.ErrorResponse
// @Failure 500 {object} utility.ErrorResponse
// @Failure default {object} utility.ErrorResponse
// @Router /api/lists/:id [get]
func GetListById(ctx utility.AppContext, service list.TodoListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		list, err := service.GetById(ctx, userId, id)
		if err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(list); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func DeleteList(ctx utility.AppContext, service list.TodoListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = service.Delete(ctx, userId, id)
		if err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(utility.StatusResponse{Status: "ok"}); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func UpdateList(ctx utility.AppContext, service list.TodoListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var input todo.UpdateListInput
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err = service.Update(ctx, userId, id, input); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(utility.StatusResponse{Status: "ok"}); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
