package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/ratelimiter/internal/pkg/jwt"
)

type Middleware struct {
}

type ctxKey string

var ownerNameKey ctxKey = "ownerName"

func New() *Middleware {
	return &Middleware{}
}

func (a *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("token")
		if token == "" {
			resp := map[string]string{
				"message": "token is required",
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(resp)
			return
		}

		claims, err := jwt.ValidateToken(token)
		if err != nil {
			resp := map[string]string{
				"message": "token is invalid",
			}
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(resp)
			return
		}

		newCtx := context.WithValue(r.Context(), ownerNameKey, claims.Ownername)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})

}

func GetOwnerName(ctx context.Context) string {
	ownerName, ok := ctx.Value(ownerNameKey).(string)
	if !ok {
		return ""
	}
	return ownerName
}
