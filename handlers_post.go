package main

import (
	"fmt"
	"net/http"

	"github.com/sugoiServal/go-rss/internal/database"
)

// get user Name from json body, then create a user with the name use database api
func (cfg *apiConfig) handlerGetPostForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't get posts %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, databasePostsToPosts(posts))
}
