package server

import (
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func getUserIDFromContext(w http.ResponseWriter, r *http.Request) (int64, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("JWT claims error: ", err)
		return 0, err
	}
	return int64(claims["id"].(float64)), nil
}
