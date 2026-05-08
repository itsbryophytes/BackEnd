package migrations

import "gorm.io/gorm"

func EnablePgVector(db *gorm.DB) error {
	return db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error
}
