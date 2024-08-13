package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"chad-rss/cmd/web"
	database "chad-rss/internal/database/sqlc"
	"chad-rss/internal/utils"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.NotFound(templ.Handler(web.NotFound()).ServeHTTP)

	// Public assets
	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", templ.Handler(web.Dashboard()).ServeHTTP)

		r.Get("/signin", templ.Handler(web.SigninForm()).ServeHTTP)
		r.Post("/signin", s.Signin)
		r.Get("/signup", templ.Handler(web.SignupForm()).ServeHTTP)
		r.Post("/signup", s.Signup)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Get("/health", s.healthHandler)
		// r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		// 	_, claims, _ := jwtauth.FromContext(r.Context())
		// 	w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		// })
	})

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Parse form error: ", err)
		return
	}
	username := r.FormValue("Username")
	password := r.FormValue("Password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Hash password error: ", err)
		return
	}

	if _, err = s.db.Query().CreateUser(context.Background(), database.CreateUserParams{Username: username, Password: string(hashedPassword)}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Create user in DB error:", err)
		return
	}

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"username":      username,
		"refresh_token": uuid.New().String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating JWT", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	})

	w.Header().Add("HX-Redirect", "/")
}

func (s *Server) Signin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	username := r.FormValue("Username")
	password := r.FormValue("Password")

	user, err := s.db.Query().GetUserByUsername(context.Background(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new token
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"refresh_token": uuid.New().String(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	})

	w.Header().Add("HX-Redirect", "/")
}

