package service

import (
	"context"
	"time"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/Caknoooo/go-gin-clean-starter/modules/profile/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/profile/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (dto.HealthProfileResponse, error)
	UpsertProfile(ctx context.Context, userID string, req dto.HealthProfileRequest) (dto.HealthProfileResponse, error)
}

type profileService struct {
	repo repository.ProfileRepository
	db   *gorm.DB
}

func NewProfileService(repo repository.ProfileRepository, db *gorm.DB) ProfileService {
	return &profileService{repo: repo, db: db}
}

func (s *profileService) GetProfile(ctx context.Context, userID string) (dto.HealthProfileResponse, error) {
	profile, err := s.repo.GetByUserID(ctx, s.db, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.HealthProfileResponse{UserID: userID}, nil
		}
		return dto.HealthProfileResponse{}, err
	}
	return toResponse(profile), nil
}

func (s *profileService) UpsertProfile(ctx context.Context, userID string, req dto.HealthProfileRequest) (dto.HealthProfileResponse, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return dto.HealthProfileResponse{}, err
	}

	var dob *time.Time
	if req.DateOfBirth != "" {
		// Try YYYY-MM-DD first
		t, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			dob = &t
		} else {
			// Try RFC3339 as fallback
			t, err = time.Parse(time.RFC3339, req.DateOfBirth)
			if err == nil {
				dob = &t
			}
		}
	}

	profile := entities.UserHealthProfile{
		UserID:             parsedUserID,
		DateOfBirth:        dob,
		BiologicalSex:      req.BiologicalSex,
		HeightCm:           req.HeightCm,
		WeightKg:           req.WeightKg,
		BloodType:          req.BloodType,
		SmokingStatus:      req.SmokingStatus,
		ExistingConditions: req.ExistingConditions,
		CurrentMedications: req.CurrentMedications,
	}

	saved, err := s.repo.Upsert(ctx, s.db, profile)
	if err != nil {
		return dto.HealthProfileResponse{}, err
	}
	return toResponse(saved), nil
}

func toResponse(p entities.UserHealthProfile) dto.HealthProfileResponse {
	var dobStr string
	if p.DateOfBirth != nil {
		dobStr = p.DateOfBirth.Format("2006-01-02")
	}

	return dto.HealthProfileResponse{
		ID:                 p.ID.String(),
		UserID:             p.UserID.String(),
		DateOfBirth:        dobStr,
		BiologicalSex:      p.BiologicalSex,
		HeightCm:           p.HeightCm,
		WeightKg:           p.WeightKg,
		BloodType:          p.BloodType,
		SmokingStatus:      p.SmokingStatus,
		ExistingConditions: p.ExistingConditions,
		CurrentMedications: p.CurrentMedications,
	}
}
