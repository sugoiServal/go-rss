package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/sugoiServal/go-rss/internal/database"
)

// get user Name from json body, then create a user with the name use database api
func (cfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feed_follow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create follow")
		return
	}

	respondWithJson(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feed_follow))
}

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := cfg.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get feeds for user: %v", user.ID))
		return
	}
	respondWithJson(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feeds))
}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	idStr := chi.URLParam(r, "followID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't parse follow id to uuid: %v", err))
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     id,
		UserID: user.ID,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't unfollow: %v", err))
		return
	}
	type msg struct {
		Msg string `json:"msg"`
	}
	respondWithJson(w, http.StatusOK, msg{
		Msg: "unfollow successful",
	})
}
