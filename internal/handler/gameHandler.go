package handler

import (
	"boarderGameStat/internal/model"
	"boarderGameStat/internal/repository"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type GameHandler struct {
	gameRepo repository.GameRepository
}

func NewGameHandler(gameRepo repository.GameRepository) *GameHandler {
	return &GameHandler{
		gameRepo: gameRepo,
	}
}

func (h GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	games, err := h.gameRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(games)
}

func (h GameHandler) GetGameByID(w http.ResponseWriter, r *http.Request) {
	gameID, err := strconv.ParseUint(r.PathValue("gameID"), 10, 64)
	if err != nil {
		log.Printf("Invalide game ID: %v", err)
		http.Error(w, "Invalide game ID: must be a number", http.StatusBadRequest)
		return
	}

	game, err := h.gameRepo.GetByID(uint(gameID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Game not found: %v", err)
			http.Error(w, "Game not found", http.StatusNotFound)
		}
		log.Printf("Faled to get game ID=%d: %v", gameID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(game)

}

func (h GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game model.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.gameRepo.Create(&game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}

func (h GameHandler) GetGameUsers(w http.ResponseWriter, r *http.Request) {
	gameID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid game id: %v", err)
		http.Error(w, "Invalid game id", http.StatusBadRequest)
		return
	}
	game, err := h.gameRepo.GetByIdWithUsers(uint(gameID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Game not found: ID=%d", gameID)
			http.Error(w, "Game not found", http.StatusNotFound)
		}
		log.Printf("Failed to get game users for ID=%d", gameID)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(game.Users); err != nil {
		log.Printf("Faled t encode game users for ID=%d: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h GameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	gameID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid game id: %v", err)
		http.Error(w, "Invalid game id", http.StatusBadRequest)
		return
	}

	if err = h.gameRepo.Delete(uint(gameID)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Faled to delete game for ID=%d", gameID)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Printf("")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}
