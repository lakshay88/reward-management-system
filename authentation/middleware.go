package auth

import (
	"net/http"
	"strings"

	utils "github.com/lakshay88/reward-management-system/Utils"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header missing"})
				return
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			// Validate token
			_, err := ValidateToken(tokenString)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized: invalid token"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
