package handler

import (
	"boarderGameStat/internal/model"
	"boarderGameStat/internal/repository"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
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

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		log.Printf("Failed to get all users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user id: must be a number", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByID(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found: ID=%d", userID)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to get user ID=%d: %v", userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Failed to decode user request: %v", err)
		http.Error(w, "Invaled request body", http.StatusBadRequest)
		return
	}
	if err := h.userRepo.Create(&user); err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) AddGameToUser(w http.ResponseWriter, r *http.Request) {
	// Парсим user ID
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		http.Error(w, "Invalid user id: must be a positive number", http.StatusBadRequest)
		return
	}

	gameID, err := strconv.ParseUint(r.PathValue("gameID"), 10, 64)
	if err != nil {
		log.Printf("Invalid game ID: %v", err)
		http.Error(w, "Invalid game id: must be a positive number", http.StatusBadRequest)
		return
	}

	if err := h.userRepo.AddGame(uint(userID), uint(gameID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User or game not found: userID=%d, gameID=%d", userID, gameID)
			http.Error(w, "User or game not found", http.StatusNotFound)
			return
		}

		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			log.Printf("Game already added to user: userID=%d, gameID=%d", userID, gameID)
			http.Error(w, "Game already added to user", http.StatusConflict)
			return
		}

		log.Printf("Failed to add game %d to user %d: %v", gameID, userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Game added to user successfully",
		"user_id": uint(userID),
		"game_id": uint(gameID),
	})
}

func (h *UserHandler) GetUserGames(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid user ID: %v, error: %v", userID, err)
		http.Error(w, "Invalid user id: must be a positive number", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByIDWithGames(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found: ID=%d", userID)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to get user games for ID=%d: %v", userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user.Games); err != nil {
		log.Printf("Failed to encode user games for ID=%d: %v", userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid user ID:  %v", err)
		http.Error(w, "Invalid user id: must be a number", http.StatusBadRequest)
		return
	}

	err = h.userRepo.Delete(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found: ID=%d", userID)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to delete user ID=%d: %v", userID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("User deleted successfully: ID=%d", userID)
	w.WriteHeader(http.StatusNoContent)
}
