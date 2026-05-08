package service

import (
	"context"

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

	profile := entities.UserHealthProfile{
		UserID:             parsedUserID,
		Age:                req.Age,
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
	return dto.HealthProfileResponse{
		ID:                 p.ID.String(),
		UserID:             p.UserID.String(),
		Age:                p.Age,
		BiologicalSex:      p.BiologicalSex,
		HeightCm:           p.HeightCm,
		WeightKg:           p.WeightKg,
		BloodType:          p.BloodType,
		SmokingStatus:      p.SmokingStatus,
		ExistingConditions: p.ExistingConditions,
		CurrentMedications: p.CurrentMedications,
	}
}
