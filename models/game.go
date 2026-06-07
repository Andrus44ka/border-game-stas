package models

import "gorm.io/gorm"

// Game представляет таблицу игр
type Game struct {
	gorm.Model
	Title        string `gorm:"type:varchar(100);not null"`
	Description  string `gorm:"type:text"`
	CountPlayers int    `gorm:"type:int;default:2;not null"`
	// CountPlayers int    `gorm:"type:int;default:2;not null"`

	// Связь с пользователями
	Users []User `gorm:"many2many:relation_user_games;"`
}

func (Game) TableName() string {
	return "games"
}
