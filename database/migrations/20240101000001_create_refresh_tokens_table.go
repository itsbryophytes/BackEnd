package migrations

import (
	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"gorm.io/gorm"
)

func Up20240101000001CreateRefreshTokensTable(db *gorm.DB) error {
	return db.AutoMigrate(&entities.RefreshToken{})
}

func Down20240101000001CreateRefreshTokensTable(db *gorm.DB) error {
	return db.Migrator().DropTable(&entities.RefreshToken{})
}
