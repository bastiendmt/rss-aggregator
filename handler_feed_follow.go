package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bastiendmt/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error while parsing JSON: %s", err))
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %s", err))
		return
	}

	respondWithJSON(w, 201, dataFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %s", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handleDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow ID: %s", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{ID: feedFollowID, UserID: user.ID})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow ID: %s", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
