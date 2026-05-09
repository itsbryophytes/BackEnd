package database

import (
	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/Caknoooo/go-gin-clean-starter/database/migrations"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(
		"20240101000000_create_users_table",
		migrations.Up20240101000000CreateUsersTable,
		migrations.Down20240101000000CreateUsersTable,
	)
	RegisterMigration(
		"20240101000001_create_refresh_tokens_table",
		migrations.Up20240101000001CreateRefreshTokensTable,
		migrations.Down20240101000001CreateRefreshTokensTable,
	)
}

func Migrate(db *gorm.DB) error {
	if err := migrations.EnablePgVector(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(
		&entities.Migration{},
		&entities.User{},
		&entities.RefreshToken{},
		&entities.UserHealthProfile{},
		&entities.File{},
		&entities.Document{},
		&entities.OCRResult{},
		&entities.ExtractedMetric{},
		&entities.HealthMetric{},
		&entities.Embedding{},
		&entities.ChatSession{},
		&entities.ChatMessage{},
		&entities.Article{},
		&entities.BloodSugarLog{},
		&entities.BloodPressureLog{},
		&entities.WeightLog{},
	); err != nil {
		return err
	}

	db.Exec(`CREATE INDEX IF NOT EXISTS embeddings_vector_idx 
             ON embeddings USING hnsw (embedding vector_cosine_ops)`)
	manager := NewMigrationManager(db)
	return manager.Run()
}
