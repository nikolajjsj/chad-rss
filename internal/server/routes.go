package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	// "github.com/go-shiori/go-readability"
	// "github.com/mmcdole/gofeed"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(Authenticator(tokenAuth))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/health", s.healthHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/signin", s.signin)
			r.Post("/signup", s.signup)
		})
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(StrictAuthenticator(tokenAuth))
		//
		// r.Get("/feeds/sidebar", s.FeedsSidebarHandler)
		//
		// r.Get("/feeds/create", templ.Handler(web.FeedsCreate()).ServeHTTP)
		// r.Post("/feeds/create", s.FeedCreateHandler)
		// r.Get("/feeds/{slug}", s.FeedHandler)
		// r.Delete("/feeds/{slug}", s.FeedDeleteHandler)
		// r.Post("/feeds/{slug}/sync", s.FeedSyncHandler)
		// r.Get("/feeds/{slug}/articles", s.FeedArticlesHandler)
		// r.Get("/feeds/{slug}/articles/{id}", s.FeedHandler)
		//
		// r.Get("/articles/{slug}", s.ArticleHandler)
		// r.Get("/articles/{slug}/content", s.ArticleContentHandler)
	})

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal("OK")
	_, _ = w.Write(jsonResp)
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Parse form error: ", err)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Decode JSON error: ", err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Hash password error: ", err)
		return
	}

	user, err := s.db.CreateUser(creds.Username, string(hashedPassword))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Create user in DB error:", err)
		return
	}

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"id":            user.ID,
		"username":      creds.Username,
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

	response, err := json.Marshal(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(response)
}

func (s *Server) signin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Decode JSON error: ", err)
		return
	}

	user, err := s.db.GetUserByUsername(creds.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
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

	response, err := json.Marshal(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(response)
}
