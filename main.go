package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	godotenv.Load(".env")
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not found in Env Variable")
	}
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

	// Define routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Chi Router!\n"))
		fmt.Fprintf(w, "%s", time.Now()) // alt way to write res
	})

	// v1-Router (/v1/)
	v1Router := chi.NewRouter()
	v1Router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		respondWithJson(w, 200, struct{}{})
	})
	v1Router.Get("/client-err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, 400, "Something went wrong")
	})
	r.Mount("/v1", v1Router)

	// Sub-router (/auth/login)
	authRouter := chi.NewRouter()
	authRouter.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("login"))
	})
	r.Mount("/auth", authRouter)

	log.Printf("Server starting on port %v", portStr)
	err := http.ListenAndServe(":"+portStr, r)
	if err != nil {
		log.Fatal(err)
	}
}
