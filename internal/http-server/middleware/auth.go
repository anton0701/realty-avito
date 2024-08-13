package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: вынести строки в константы
// TODO: подумать над UserType

var jwtKey = []byte("jwt_most_secret_key")

type Claims struct {
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func GenerateJWT(userType string) (string, error) {
	claims := &Claims{
		UserType:         userType,
		RegisteredClaims: jwt.RegisteredClaims{ID: generateUUID()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_type", claims.UserType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireModerator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType := r.Context().Value("user_type").(string)
		if userType != "moderator" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
