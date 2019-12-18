package middleware

import (
	"context"
	"net/http"
	"os"

	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/dgrijalva/jwt-go"
)

type JwtContextKey string

type Token struct {
	Email        string `json:"email"`
	Role         string `json:"role"`
	Organization string `json:"organization"`
	jwt.StandardClaims
}


// JwtAuthentication middleware for authorizing endpoints
func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//notAuth := []string{"/api/v1/admin/ping", "/api/v1/participants/save-attendee"}   // List of endpoints that doesn't require auth
		//requestPath := r.URL.Path // Current request path

		// Check if request does not need authentication, serve the request if it doesn't need it
		//for _, value := range notAuth {
		//	if value == requestPath {
		//		next.ServeHTTP(w, r)
		//		return
		//	}
		//}

		var response map[string]interface{}
		tokenHeader := r.Header.Get("Authorization") // Grab the token from the header

		if tokenHeader == "" { // Token is missing, returns with error code 403 Unauthorized
			response = u.Message(http.StatusForbidden, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})

		if err != nil { // Malformed token, returns with http code 403 as usual
			response = u.Message(http.StatusForbidden, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { // Token is invalid, maybe not signed on this server
			response = u.Message(http.StatusForbidden, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		// Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), JwtContextKey("token"), tk)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) // proceed in the middleware chain
	})
}
