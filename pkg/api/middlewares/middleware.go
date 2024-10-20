package middlewares

import (
	"context"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/auth"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	UserCtx             = "UserId"
)

type UserAuthMiddleware struct {
	service auth.AuthorizationService
}

func NewUserAuthMiddleware(service auth.AuthorizationService) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		service: service,
	}
}

func (m *UserAuthMiddleware) UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			utility.NewErrorResponse(w, http.StatusUnauthorized, "empty auth header")
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			utility.NewErrorResponse(w, http.StatusUnauthorized, "invalid auth header")
			return
		}

		userId, err := m.service.ParseToken(headerParts[1])
		if err != nil {
			utility.NewErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserCtx, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(ctx utility.AppContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Logger.Error("Panic occurred",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())))
				utility.NewErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
			next.ServeHTTP(w, r)
		}()
	})
}
