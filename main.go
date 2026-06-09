package main

import (
	"log"
	"net/http"

	"boarderGameStat/internal/database"
	"boarderGameStat/internal/handler"
	"boarderGameStat/internal/repository"
)

func main() {
	db := database.New()

	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)

	userHandler := handler.NewUserHandler(*gameRepo, *userRepo)
	gameHandler := handler.NewGameHandler(*gameRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", userHandler.GetUser)
	mux.HandleFunc("GET /users/{id}/games", userHandler.GetUserByID)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("POST /users/{id}/games/{gameID}", userHandler.AddGameToUser)

	mux.HandleFunc("GET /games", gameHandler.GetGame)
	mux.HandleFunc("POST /games", gameHandler.CreateGame)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
