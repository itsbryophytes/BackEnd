package migrations

import (
	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"gorm.io/gorm"
)

func Up20240101000000CreateUsersTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.User{})
}

func Down20240101000000CreateUsersTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.User{})
}
