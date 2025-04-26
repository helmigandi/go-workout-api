package api

import (
	"encoding/json"
	"github.com/helmigandi/go-workout-api/internal/store"
	"github.com/helmigandi/go-workout-api/internal/tokens"
	"github.com/helmigandi/go-workout-api/internal/utils"
	"log"
	"net/http"
	"time"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (t *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var request createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		t.logger.Printf("ERROR: decoding create token request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request send"})
		return
	}

	user, err := t.userStore.GetUserByUsername(request.Username)
	if err != nil || user == nil {
		t.logger.Printf("ERROR: GetUserByUsername: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	matches, err := user.PasswordHash.Matches(request.Password)
	if err != nil {
		t.logger.Printf("ERROR: PasswordHash.Matches: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !matches {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid username or password"})
		return
	}

	token, err := t.tokenStore.CreateToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		t.logger.Printf("ERROR: CreateToken: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
