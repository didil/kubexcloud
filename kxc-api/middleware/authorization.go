package middleware

import (
	"net/http"

	"github.com/didil/kubexcloud/kxc-api/handlers"
)

// Authorization middleware
func Authorization(root *handlers.Root, role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userName := r.Context().Value(handlers.CtxKey("userName")).(string)

			ok, err := root.UserSvc.HasRole(r.Context(), userName, role)
			if err != nil {
				root.HandleError(w, r, err)
				return
			}

			if !ok {
				handlers.JSONError(w, "not authorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
