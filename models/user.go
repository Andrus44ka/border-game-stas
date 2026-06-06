package models

import "gorm.io/gorm"

// User представляет таблицу пользователей
type User struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100)"`
	Login        string `gorm:"type:varchar(50);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	FullName     string `gorm:"type:varchar(100)"`
}

func (User) TableName() string {
	return "users"
}
