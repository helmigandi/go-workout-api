package api

import (
	"encoding/json"
	"errors"
	"github.com/helmigandi/go-workout-api/internal/store"
	"github.com/helmigandi/go-workout-api/internal/utils"
	"log"
	"net/http"
	"regexp"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (u *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("username must be less than 50 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	return nil
}

func (u *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var request registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		u.logger.Printf("ERROR: decoding register user request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request send"})
		return
	}

	err = u.validateRegisterRequest(&request)
	if err != nil {
		u.logger.Printf("ERROR: validate register user request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: request.Username,
		Email:    request.Email,
	}

	if request.Bio != "" {
		user.Bio = request.Bio
	}

	err = user.PasswordHash.Set(request.Password)
	if err != nil {
		u.logger.Printf("ERROR: set password hash: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = u.userStore.CreateUser(user)
	if err != nil {
		u.logger.Printf("ERROR: register user: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}
