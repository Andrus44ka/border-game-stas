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

	mux.HandleFunc("GET /users", userHandler.GetAllUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("GET /users/{id}/games", userHandler.GetUserGames)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

	mux.HandleFunc("POST /users/{id}/games/{gameID}", userHandler.AddGameToUser)

	mux.HandleFunc("GET /games", gameHandler.GetGame)
	mux.HandleFunc("GET /games/{id}", gameHandler.GetGameByID)
	mux.HandleFunc("GET /game/{id}/users", gameHandler.GetGameUsers)
	mux.HandleFunc("POST /games", gameHandler.CreateGame)
	mux.HandleFunc("DELETE /game/{id}", gameHandler.DeleteGame)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
