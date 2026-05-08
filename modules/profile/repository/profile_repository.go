package repository

import (
	"context"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProfileRepository interface {
	GetByUserID(ctx context.Context, db *gorm.DB, userID string) (entities.UserHealthProfile, error)
	Upsert(ctx context.Context, db *gorm.DB, profile entities.UserHealthProfile) (entities.UserHealthProfile, error)
}

type profileRepository struct{}

func NewProfileRepository() ProfileRepository {
	return &profileRepository{}
}

func (r *profileRepository) GetByUserID(ctx context.Context, db *gorm.DB, userID string) (entities.UserHealthProfile, error) {
	var profile entities.UserHealthProfile
	err := db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	return profile, err
}

func (r *profileRepository) Upsert(ctx context.Context, db *gorm.DB, profile entities.UserHealthProfile) (entities.UserHealthProfile, error) {
	err := db.WithContext(ctx).
		Where(entities.UserHealthProfile{UserID: profile.UserID}).
		Assign(profile).
		FirstOrCreate(&profile).Error
	if err != nil {
		return entities.UserHealthProfile{}, err
	}
	// If record already existed, update all fields
	err = db.WithContext(ctx).
		Model(&profile).
		Where("user_id = ?", profile.UserID).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"age", "biological_sex", "height_cm", "weight_kg", "blood_type", "smoking_status", "existing_conditions", "current_medications", "updated_at"}),
		}).
		Save(&profile).Error
	return profile, err
}
