package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"boarderGameStat/internal/database"
	"boarderGameStat/internal/model"
	"boarderGameStat/internal/repository"
)

func main() {
	db := database.New()

	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)

	http.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		users, err := userRepo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	})

	http.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := userRepo.Create(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	})

	http.HandleFunc("GET /users/{id}/games", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		user, err := userRepo.GetByIDWithGames(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	})

	http.HandleFunc("POST /users/{id}/games/{gameID}", func(w http.ResponseWriter, r *http.Request) {
		userID, _ := strconv.Atoi(r.PathValue("id"))
		gameID, _ := strconv.Atoi(r.PathValue("gameID"))
		if err := userRepo.AddGame(uint(userID), uint(gameID)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("GET /games", func(w http.ResponseWriter, r *http.Request) {
		games, err := gameRepo.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(games)
	})

	http.HandleFunc("POST /games", func(w http.ResponseWriter, r *http.Request) {
		var game model.Game
		if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := gameRepo.Create(&game); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(game)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
