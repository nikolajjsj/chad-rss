package server

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				ctx := context.WithValue(r.Context(), "authenticated", false)
				next.ServeHTTP(w, r.WithContext(ctx))
			}

			if token == nil || jwt.Validate(token, tokenAuth.ValidateOptions()...) != nil {
				ctx := context.WithValue(r.Context(), "authenticated", false)
				next.ServeHTTP(w, r.WithContext(ctx))
			}

			// Token is authenticated, pass it through
			ctx := context.WithValue(r.Context(), "authenticated", true)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

// func SetupContext(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.WithValue(r.Context(), "authenticated", false)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
