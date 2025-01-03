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
	"strconv"
	"strings"
	"time"

	"chad-rss/cmd/web"
	database "chad-rss/internal/database/sqlc"
	"chad-rss/internal/utils"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-shiori/go-readability"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mmcdole/gofeed"
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
	r.Use(middleware.Recoverer)
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(Authenticator(tokenAuth))

	// 404 Handler
	r.NotFound(templ.Handler(web.NotFound()).ServeHTTP)

	// Public assets
	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/signin", templ.Handler(web.SigninForm()).ServeHTTP)
		r.Post("/signin", s.Signin)
		r.Get("/signup", templ.Handler(web.SignupForm()).ServeHTTP)
		r.Post("/signup", s.Signup)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(StrictAuthenticator(tokenAuth))

		r.Get("/health", s.healthHandler)

		r.Get("/", templ.Handler(web.Dashboard()).ServeHTTP)

		r.Get("/feeds/sidebar", s.FeedsSidebarHandler)

		r.Get("/feeds/create", templ.Handler(web.FeedsCreate()).ServeHTTP)
		r.Post("/feeds/create", s.FeedCreateHandler)
		r.Get("/feeds/{slug}", s.FeedHandler)
		r.Delete("/feeds/{slug}", s.FeedDeleteHandler)
		r.Post("/feeds/{slug}/sync", s.FeedSyncHandler)
		r.Get("/feeds/{slug}/articles", s.FeedArticlesHandler)
		r.Get("/feeds/{slug}/articles/{id}", s.FeedHandler)

		r.Get("/articles/{slug}", s.ArticleHandler)
		r.Get("/articles/{slug}/content", s.ArticleContentHandler)
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

	user, err := s.db.Query().CreateUser(context.Background(), database.CreateUserParams{Username: username, Password: string(hashedPassword)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Create user in DB error:", err)
		return
	}

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"id":            user.ID,
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
	log.Println(tokenString)
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

func (s *Server) FeedsSidebarHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(w, r)
	if err != nil {
		return
	}

	feeds, err := s.db.Query().GetFeeds(context.Background(), database.GetFeedsParams{
		ID:     userID,
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	web.SidebarFeeds(feeds).Render(context.Background(), w)
}

func (s *Server) FeedHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(w, r)
	if err != nil {
		return
	}
	feedID := chi.URLParam(r, "slug")
	articleID := chi.URLParam(r, "id")

	feed, err := s.db.Query().GetFeedByID(context.Background(), database.GetFeedByIDParams{
		Nid: feedID,
		ID:  userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if articleID != "" {
		web.Feed(web.FeedProps{
			Feed:      feed,
			ArticleID: articleID,
		}).Render(context.Background(), w)
	} else {
		web.Feed(web.FeedProps{
			Feed: feed,
		}).Render(context.Background(), w)
	}

}

func (s *Server) FeedArticlesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(w, r)
	if err != nil {
		return
	}
	feedID := chi.URLParam(r, "slug")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	var PER_PAGE int64 = 20
	articles, err := s.db.Query().GetUserFeedArticles(context.Background(), database.GetUserFeedArticlesParams{
		ID:     userID,
		Nid:    feedID,
		Limit:  PER_PAGE,
		Offset: (int64(page) - 1) * PER_PAGE,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	web.FeedList(feedID, page+1, articles).Render(context.Background(), w)
}

func (s *Server) FeedCreateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Parse form error: ", err)
		return
	}
	feedURL := r.FormValue("URL")

	_, err := url.ParseRequestURI(feedURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid URL: ", err)
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Parse feed error: ", err)
		return
	}

	userID, err := getUserIDFromContext(w, r)
	if err != nil {
		return
	}
	feedID, err := utils.GenerateNID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Generate NID error: ", err)
		return
	}

	ctx := context.Background()
	tx, err := s.db.Transaction(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("DB transaction error: ", err)
		return
	}
	defer tx.Rollback()

	qtx := s.db.Query().WithTx(tx)

	image := ""
	if feed.Image != nil && feed.Image.URL != "" {
		image = feed.Image.URL
	}

	feedDB, err := qtx.CreateFeed(ctx, database.CreateFeedParams{
		Nid:     feedID,
		Url:     feedURL,
		Title:   feed.Title,
		Summary: sql.NullString{String: feed.Description, Valid: feed.Description != ""},
		Image:   sql.NullString{String: image, Valid: feed.Image != nil && feed.Image.URL != ""},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Create feed in DB error: ", err)
		return
	}

	if err := qtx.AddFeedToUser(ctx, database.AddFeedToUserParams{
		UserID: userID,
		FeedID: feedDB.ID,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Add feed to user in DB error: ", err)
		return
	}

	tx.Commit()

	w.Header().Add("HX-Redirect", fmt.Sprintf("/feeds/%s", feedID))
}

func (s *Server) FeedSyncHandler(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "slug")

	userID, err := getUserIDFromContext(w, r)
	if err != nil {
		return
	}

	feedDB, err := s.db.Query().GetFeedByID(context.Background(), database.GetFeedByIDParams{
		Nid: feedID,
		ID:  userID,
	})

	fp := gofeed.NewParser()
	f, err := fp.ParseURL(feedDB.Url)
	if err != nil {
		log.Println("Parse feed error: ", err)
		return
	}

	for _, item := range f.Items {
		// TODO: move this to a separate function
		if item.GUID == "" || item.Link == "" || item.Title == "" {
			continue
		}

		media := ""
		if item.Image != nil && item.Image.URL != "" {
			media = item.Image.URL
		} else {
			for i := range item.Links {
				if strings.Contains(item.Links[i], ".jpeg") || strings.Contains(item.Links[i], ".jpg") || strings.Contains(item.Links[i], ".png") {
					media = item.Links[i]
					break
				}
			}
		}

		authors := ""
		if item.Authors != nil {
			stringSlice := make([]string, len(item.Authors))
			for i, v := range item.Authors {
				stringSlice[i] = v.Name
			}
			authors = strings.Join(stringSlice[:], ",")
		}

		nid, err := utils.GenerateNID()
		if err != nil {
			log.Println("Generate NID error: ", err)
			continue
		}

		if _, err = s.db.Query().CreateFeedArticles(context.Background(), database.CreateFeedArticlesParams{
			Nid:         nid,
			RssID:       item.GUID,
			Url:         item.Link,
			Title:       item.Title,
			Summary:     sql.NullString{String: item.Description, Valid: item.Description != ""},
			Content:     sql.NullString{String: item.Content, Valid: true},
			Authors:     sql.NullString{String: authors, Valid: authors != ""},
			Media:       sql.NullString{String: media, Valid: media != ""},
			PublishedAt: sql.NullTime{Time: *item.PublishedParsed, Valid: item.PublishedParsed != nil},
			FeedID:      feedDB.ID,
		}); err != nil {
			log.Println("Create article in DB error: ", err)
			continue
		}
	}
}

func (s *Server) FeedDeleteHandler(w http.ResponseWriter, r *http.Request) {
	feedID := chi.URLParam(r, "slug")
	userID, err := getUserIDFromContext(w, r)
	if err != nil {
		return
	}

	feed, err := s.db.Query().GetFeedByID(context.Background(), database.GetFeedByIDParams{
		Nid: feedID,
		ID:  userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Get feed by ID in DB error: ", err)
		return
	}

	if err = s.db.Query().RemoveFeedFromUser(context.Background(), database.RemoveFeedFromUserParams{
		UserID: userID,
		FeedID: feed.ID,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Remove feed from user in DB error: ", err)
		return
	}

	w.Header().Add("HX-Redirect", "/")
}

func (s *Server) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "slug")
	userID, err := getUserIDFromContext(w, r)
	if err != nil {
		return
	}

	article, err := s.db.Query().GetArticle(context.Background(), database.GetArticleParams{ID: userID, Nid: articleID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	web.Article(article).Render(context.Background(), w)
}

func (s *Server) ArticleContentHandler(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "slug")
	userID, err := getUserIDFromContext(w, r)
	if err != nil {
		return
	}

	article, err := s.db.Query().GetArticle(context.Background(), database.GetArticleParams{ID: userID, Nid: articleID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	contentType := r.URL.Query().Get("content")

	if contentType == "extracted" {
		url := article.Url
		content, err := readability.FromURL(url, 5*time.Second)
		if err != nil {
			_, _ = w.Write([]byte("Not able to extract content from this article"))
		}

		_, _ = w.Write([]byte(content.Content))
	} else {
		_, _ = w.Write([]byte(article.Content.String))
	}
}
