package api

import (
	"context"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/handler"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/middlewares"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/auth"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/item"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/list"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"
)

const (
	maxHeaderBytes = 1 << 20 // 1 MB
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	subRouter  *mux.Router
}

func NewServer(ctx utility.AppContext, middleware *middlewares.UserAuthMiddleware) *Server {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.UserAuth)

	wrappedRouter := middlewares.RecoveryMiddleware(ctx, router)

	return &Server{
		httpServer: &http.Server{
			Addr:           ":8000",
			MaxHeaderBytes: maxHeaderBytes,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			Handler:        wrappedRouter,
		},
		router:    router,
		subRouter: api,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) HandleAuth(service auth.AuthorizationService) {
	s.router.HandleFunc("/auth/sign-up/", handler.SignUp(service)).Methods(http.MethodPost)
	s.router.HandleFunc("/auth/sign-in/", handler.SignIn(service)).Methods(http.MethodPost)
}

func (s *Server) HandleLists(ctx utility.AppContext, service list.TodoListService) {
	s.subRouter.HandleFunc("/lists/", handler.CreateList(ctx, service)).Methods(http.MethodPost)
	s.subRouter.HandleFunc("/lists/", handler.GetAllLists(ctx, service)).Methods(http.MethodGet)
	s.subRouter.HandleFunc("/lists/{id}", handler.GetListById(ctx, service)).Methods(http.MethodGet)
	s.subRouter.HandleFunc("/lists/{id}", handler.DeleteList(ctx, service)).Methods(http.MethodDelete)
	s.subRouter.HandleFunc("/lists/{id}", handler.UpdateList(ctx, service)).Methods(http.MethodPut)
}

func (s *Server) HandleItems(ctx utility.AppContext, service item.TodoItemService) {
	s.subRouter.HandleFunc("/lists/{id}/items/", handler.CreateItem(ctx, service)).Methods(http.MethodPost)
	s.subRouter.HandleFunc("/lists/{id}/items/", handler.GetAllItems(ctx, service)).Methods(http.MethodGet)
	s.subRouter.HandleFunc("/items/{id}", handler.GetItemById(ctx, service)).Methods(http.MethodGet)
	s.subRouter.HandleFunc("/items/{id}", handler.DeleteItem(ctx, service)).Methods(http.MethodDelete)
	s.subRouter.HandleFunc("/items/{id}", handler.UpdateItem(ctx, service)).Methods(http.MethodPut)
}
