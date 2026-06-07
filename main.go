package main

import (
	"boarderGameStat/models"
	"fmt"
	"log"
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
	err = models.SetupAssociations(db)
	if err != nil {
		log.Fatal("Ошибка настройки связей:", err)
	}

	fmt.Println("База данных и связи успешно настроены!")
}
