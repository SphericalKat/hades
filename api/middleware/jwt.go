package middleware

import (
	"context"
	"net/http"
	"os"

	u "github.com/ATechnoHazard/hades-2/internal/utils"
	"github.com/ATechnoHazard/janus"
	"github.com/dgrijalva/jwt-go"
)

type JwtContextKey string

type Token struct {
	Email string `json:"email"`
	OrgID uint   `json:"org_id"`
	jwt.StandardClaims
}

// JwtAuthentication middleware for authorizing endpoints
func JwtAuthentication(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response map[string]interface{}
		tokenHeader := r.Header.Get("Authorization") // Grab the token from the header

		if tokenHeader == "" { // Token is missing, returns with error code 403 Unauthorized
			response = u.Message(http.StatusForbidden, "Missing auth token")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})

		if err != nil { // Malformed token, returns with http code 403 as usual
			response = u.Message(http.StatusForbidden, "Malformed authentication token")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		if !token.Valid { // Token is invalid, maybe not signed on this server
			response = u.Message(http.StatusForbidden, "Token is not valid.")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		// Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), JwtContextKey("token"), tk)
		ctx = context.WithValue(ctx, "janus_context", &janus.Account{
			CacheKey:       tk.Email,
			OrganizationID: tk.OrgID,
			Role:           "",
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) // proceed in the middleware chain
	}
}
