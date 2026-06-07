package main

import (
	"boarderGameStat/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	// defer db.Close()

	err = db.AutoMigrate(&models.User{}, &models.Game{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	// 3. Настройка связей (из нашего нового файла)
	// err = models.SetupAssociations(db)
	// if err != nil {
	// 	log.Fatal("Ошибка настройки связей:", err)
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go + PostgreSQL!")
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
