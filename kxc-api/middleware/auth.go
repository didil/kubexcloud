package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/didil/kubexcloud/kxc-api/handlers"
	"github.com/didil/kubexcloud/kxc-api/services"
)

// Authentication middleware builder
func Authentication(root *handlers.Root) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			var userName string
			if len(authHeader) > 0 {
				authFields := strings.Fields(authHeader)
				if len(authFields) != 2 || authFields[0] != "Bearer" {
					handlers.JSONError(w, "invalid authorization", http.StatusUnauthorized)
					return
				}

				token := authFields[1]
				var err error
				userName, err = services.ParseJWT(token)
				if err != nil {
					handlers.JSONError(w, fmt.Sprintf("invalid auth token"), http.StatusUnauthorized)
					return
				}
			}

			if userName == "" {
				// no auth, fail
				handlers.JSONError(w, "no authorization", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), handlers.CtxKey("userName"), userName)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
