package middleware

import (
	"admin-template/internal/core/user/utils"
	"context"
	"net/http"
	"strings"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			//logic
			authheader := request.Header.Get("Authorization")
			if authheader == "" {
				http.Error(writer, "missing header", http.StatusUnauthorized)
				return
			}
			if !strings.HasPrefix(authheader, "Bearer") {
				http.Error(writer, "missing bearier", http.StatusUnauthorized)
				return
			}
			tokenstr := strings.TrimPrefix(authheader, "Bearer ")
			if tokenstr == "" {
				http.Error(writer, "missing token", http.StatusUnauthorized)
				return
			}
			userID, err := utils.ValidateToken(tokenstr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(request.Context(), UserIDKey, *userID)
			next.ServeHTTP(writer, request.WithContext(ctx))

		})

	}
}
