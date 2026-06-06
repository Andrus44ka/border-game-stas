package main

import (
	"boardgames/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const host = "localhost"
const user = "cergo"
const password = ""
const dbname = "boarder_games_stat_db"
const port = 5432
const sslmode = "disable"

func main() {
	// 1. Подключение
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}

	// 2. Миграция таблиц
	err = db.AutoMigrate(&models.User{}, &models.Game{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	// 3. Настройка связей (из нашего нового файла)
	err = models.SetupAssociations(db)
	if err != nil {
		log.Fatal("Ошибка настройки связей:", err)
	}

	fmt.Println("База данных и связи успешно настроены!")
}
