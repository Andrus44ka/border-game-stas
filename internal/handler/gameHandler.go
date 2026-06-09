package handler

import (
	"boarderGameStat/internal/model"
	"boarderGameStat/internal/repository"
	"encoding/json"
	"net/http"
)

type GameHandler struct {
	gameRepo repository.GameRepository
}

func NewGameHandler(gameRepo repository.GameRepository) *GameHandler {
	return &GameHandler{
		gameRepo: gameRepo,
	}
}

func (g GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	games, err := g.gameRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(games)
}

func (g GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game model.Game
	if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := g.gameRepo.Create(&game); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)
}
