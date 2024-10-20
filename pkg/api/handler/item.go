package handler

import (
	"encoding/json"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/item"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateItem(ctx utility.AppContext, service item.TodoItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		listId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var input todo.TodoItem
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		id, err := service.Create(ctx, userId, listId, input)
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

func GetAllItems(ctx utility.AppContext, service item.TodoItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		listId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		items, err := service.GetAll(ctx, userId, listId)
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(items); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func GetItemById(ctx utility.AppContext, service item.TodoItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		itemId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		item, err := service.GetById(ctx, userId, itemId)
		if err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(item); err != nil {
			utility.NewErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func DeleteItem(ctx utility.AppContext, service item.TodoItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		itemId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = service.Delete(ctx, userId, itemId)
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

func UpdateItem(ctx utility.AppContext, service item.TodoItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("UserId").(int)

		itemId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var input todo.UpdateItemInput
		if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
			utility.NewErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err = service.Update(ctx, userId, itemId, input); err != nil {
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
