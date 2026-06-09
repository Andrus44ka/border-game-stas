package handler

import (
	"boarderGameStat/internal/model"
	"boarderGameStat/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserHandler struct {
	gameRepo repository.GameRepository
	userRepo repository.UserRepository
}

func NewUserHandler(gameRepo repository.GameRepository, userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		gameRepo: gameRepo,
		userRepo: userRepo,
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.userRepo.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) AddGameToUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.PathValue("id"))
	gameID, _ := strconv.Atoi(r.PathValue("gameID"))
	if err := h.userRepo.AddGame(uint(userID), uint(gameID)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	user, err := h.userRepo.GetByIDWithGames(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
