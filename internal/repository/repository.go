package repository

import (
	"gorm.io/gorm"
	"main/internal/service"
	"main/pkg"
)

type Repository struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) service.Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByKcId(kcId string) (*pkg.User, error) {
	var user pkg.User
	if err := r.db.Select("user_id, kc_id").Where("kc_id = ?", kcId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUserWithBaseRole(user *pkg.User) error {
	tx := r.db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	u2r := pkg.User2Role{
		RoleId: 1,
		UserId: user.Id,
	}

	if err := tx.Create(&u2r).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
