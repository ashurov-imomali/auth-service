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

func (r *Repository) CreateUserWithBaseRole(user *pkg.User) error {
	tx := r.db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	u2r := pkg.User2Role{
		RoleId: 16,
		UserId: user.Id,
	}

	if err := tx.Create(&u2r).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) GetUserByKcId(kcId string) (*pkg.User, bool, error) {
	var user pkg.User
	err := r.db.Where("kc_id = ?", kcId).First(&user).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, false, nil
	}
	return &user, true, err
}

func (r *Repository) GetPermissionsByUserId(id int64) ([]int64, error) {
	permissions := make([]int64, 0)
	err := r.db.Table("trole2permissions p").Select("p.id").
		Joins("inner join tuser2role t on p.role_id = t.role_id").
		Where("p.is_active = ? and t.user_id = ?", true, id).
		Find(&permissions).Error
	return permissions, err
}

func (r *Repository) GetUserInfoByKcId(kcId string) (*pkg.UserInfo, error) {
	var user pkg.UserInfo
	return &user, r.db.Where("u.kc_id=?", kcId).Select("u.user_id, u.kc_id, u.username, r.role_name as role").
		Table("tusers u").Joins("join tuser2role ur on ur.user_id = u.user_id").
		Joins("join troles r on r.role_id = ur.role_id").First(&user).Error
}

func (r *Repository) GetUserById(id int64) (*pkg.User, error) {
	var user pkg.User
	return &user, r.db.Where("user_id=?", id).First(&user).Error
}
