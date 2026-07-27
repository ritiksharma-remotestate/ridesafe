package middleware

import (
	"context"
	"fmt"
	"net/http"

	"ridesafe/models"
	"ridesafe/repository"
	"ridesafe/utils"

	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKeys string

const (
	userContext ContextKeys = "__userContext"
)

type Claims struct {
	UserID string      `json:"user_id"`
	Email  string      `json:"email"`
	Role   models.Role `json:"role"`

	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {

			utils.RespondError(w, http.StatusUnauthorized, nil, "missing auth header")
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {

			utils.RespondError(w, http.StatusForbidden, nil, "invalid authorization header")
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {

			utils.RespondError(w, http.StatusUnauthorized, err, "invalid token")
			return
		}
		user, err := repository.GetUserByID(claims.UserID)
		if err != nil || user == nil {

			utils.RespondError(w, http.StatusNotFound, err, "user not found")
			return
		}

		ctx := context.WithValue(r.Context(), userContext, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsContext(r *http.Request) *utils.Claims {
	claims, ok := r.Context().Value(userContext).(*utils.Claims)
	if !ok {
		return nil
	}
	return claims
}

func ShouldHaveRole(next http.Handler, roles ...models.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims := ClaimsContext(r)
		if claims == nil {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
			return
		}

		for _, role := range roles {
			fmt.Println(claims.Role, role)
			if claims.Role == role {
				next.ServeHTTP(w, r)
				return
			}
		}

		utils.RespondError(w, http.StatusForbidden, nil, "forbidden")
	})

}
