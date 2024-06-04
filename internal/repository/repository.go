package repository

import (
	"errors"
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

func (r *Repository) createUserWithBaseRole(user *pkg.User) error {
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

func (r *Repository) SaveUser(kcId string) error {
	var user pkg.User
	err := r.db.Where("kc_id = ?", kcId).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user.KcId = kcId
		return r.createUserWithBaseRole(&user)
	}
	return err
}

func (r *Repository) GetUserByKcId(kcId string) (*pkg.User, error) {
	var user pkg.User
	err := r.db.Where("kc_id = ?", kcId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
