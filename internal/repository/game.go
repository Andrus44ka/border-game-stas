package repository

import (
	"boarderGameStat/internal/model"

	"gorm.io/gorm"
)

type GameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

func (r *GameRepository) Create(game *model.Game) error {
	return r.db.Create(game).Error
}

func (r *GameRepository) GetByID(id uint) (*model.Game, error) {
	var game model.Game
	err := r.db.First(&game, id).Error
	return &game, err
}

func (r *GameRepository) GetByIdWithUsers(id uint) (*model.Game, error) {
	var game model.Game
	err := r.db.Preload("Users").First(&game, id).Error
	return &game, err
}

func (r *GameRepository) GetAll() ([]model.Game, error) {
	var games []model.Game
	err := r.db.Find(&games).Error
	return games, err
}

func (r *GameRepository) GetWithPlayers(id uint) (*model.Game, error) {
	var game model.Game
	err := r.db.Preload("Users").First(&game, id).Error
	return &game, err
}

func (r *GameRepository) Update(game *model.Game) error {
	return r.db.Save(game).Error
}

func (r *GameRepository) Delete(id uint) error {
	return r.db.Delete(&model.Game{}, id).Error
}
