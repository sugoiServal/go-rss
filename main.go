package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sugoiServal/go-rss/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// get env vars
	godotenv.Load(".env")
	portStr := os.Getenv("PORT")
	dbStr := os.Getenv("DB_URL")
	if portStr == "" {
		log.Fatal("PORT is not found in Env Variable")
	}
	if dbStr == "" {
		log.Fatal("DB_URL is not found in Env Variable")
	}

	// connect to db
	db, err := sql.Open("postgres", dbStr)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	log.Println("connected to db")
	queries := database.New(db)
	apiCfg := apiConfig{
		DB: queries,
	}
	go startScraping(queries, 10, 30*time.Minute)
	// create main router: r
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// v1-Router (/v1/)
	v1Router := chi.NewRouter()

	// health check
	v1Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondWithJson(w, 200, struct{}{})
	})
	v1Router.Get("/client-err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, 400, "Something went wrong")
	})

	// user
	v1Router.Post("/users", apiCfg.handlerUsersCreate)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerUserGet))

	// feed
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedCreate))
	v1Router.Get("/feeds", apiCfg.handlerFeedGet)

	// follow feed
	v1Router.Post("/feed-follow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed-follow", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet))
	v1Router.Delete("/feed-follow/{followID}", apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsDelete))

	// post
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostForUser))

	r.Mount("/v1", v1Router)

	log.Printf("Server starting on port %v", portStr)
	http.ListenAndServe(":"+portStr, r)

}
