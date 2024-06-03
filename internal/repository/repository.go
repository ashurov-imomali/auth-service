package repository

import (
	"gorm.io/gorm"
	"main/internal/service"
	"main/pkg"
)

type Repository struct {
	db *gorm.DB
}

func IntRepository(db *gorm.DB) service.Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser() {

}

func (r *Repository) GetUserByKcId(kcId string) (*pkg.User, error) {
	var user pkg.User
	if err := r.db.Select("user_id, kc_id").Where("kc_id = ?", kcId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
