package infrastructure

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/pkg"
)

func GetDbConnection(conf *pkg.Config) (*gorm.DB, error) {
	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		conf.Db.Host, conf.Db.Port, conf.Db.User, conf.Db.Password, conf.Db.Dbname)
	return gorm.Open(postgres.Open(dbURI), &gorm.Config{})
}
