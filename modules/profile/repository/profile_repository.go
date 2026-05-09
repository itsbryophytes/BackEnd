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
	// Use Clause OnConflict to handle both Create and Update in one atomic operation
	err := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"date_of_birth", "biological_sex", "height_cm", "weight_kg", "blood_type", "smoking_status", "existing_conditions", "current_medications", "updated_at"}),
	}).Create(&profile).Error

	if err != nil {
		return entities.UserHealthProfile{}, err
	}

	// Fetch the final state to return (ensuring ID is populated if it was an insert)
	var finalProfile entities.UserHealthProfile
	err = db.WithContext(ctx).Where("user_id = ?", profile.UserID).First(&finalProfile).Error
	return finalProfile, err
}
