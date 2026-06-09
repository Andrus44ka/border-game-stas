package repository

import (
	"boarderGameStat/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetByIDWithGames(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Games").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) AddGame(userID, gameID uint) error {
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}
	var game model.Game
	if err := r.db.First(&game, gameID).Error; err != nil {
		return err
	}
	return r.db.Model(&user).Association("Games").Append(&game)
}

func (r *UserRepository) RemoveGame(userID, gameID uint) error {
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}
	var game model.Game
	if err := r.db.First(&game, gameID).Error; err != nil {
		return err
	}
	return r.db.Model(&user).Association("Games").Delete(&game)
}
