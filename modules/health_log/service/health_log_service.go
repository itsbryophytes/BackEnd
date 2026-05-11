package service

import (
	"context"
	"math"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/Caknoooo/go-gin-clean-starter/modules/health_log/dto"
	"github.com/Caknoooo/go-gin-clean-starter/modules/health_log/repository"
	"github.com/google/uuid"
)

type HealthLogService interface {
	CreateBloodSugar(ctx context.Context, userID uuid.UUID, req dto.CreateBloodSugarRequest) (entities.BloodSugarLog, error)
	GetBloodSugarLogs(ctx context.Context, userID uuid.UUID) ([]entities.BloodSugarLog, error)
	DeleteBloodSugar(ctx context.Context, userID uuid.UUID, id uuid.UUID) error

	CreateBloodPressure(ctx context.Context, userID uuid.UUID, req dto.CreateBloodPressureRequest) (entities.BloodPressureLog, error)
	GetBloodPressureLogs(ctx context.Context, userID uuid.UUID) ([]entities.BloodPressureLog, error)
	DeleteBloodPressure(ctx context.Context, userID uuid.UUID, id uuid.UUID) error

	CreateWeight(ctx context.Context, userID uuid.UUID, req dto.CreateWeightRequest) (entities.WeightLog, error)
	GetWeightLogs(ctx context.Context, userID uuid.UUID) ([]entities.WeightLog, error)
	DeleteWeight(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type healthLogService struct {
	repo repository.HealthLogRepository
}

func NewHealthLogService(repo repository.HealthLogRepository) HealthLogService {
	return &healthLogService{repo: repo}
}

func (s *healthLogService) CreateBloodSugar(ctx context.Context, userID uuid.UUID, req dto.CreateBloodSugarRequest) (entities.BloodSugarLog, error) {
	indicator := "normal"
	if req.GlucoseValue < 70 {
		indicator = "low"
	} else if req.GlucoseValue > 140 {
		indicator = "high"
	}

	log := entities.BloodSugarLog{
		UserID:          userID,
		RecordedAt:      req.RecordedAt,
		GlucoseValue:    req.GlucoseValue,
		MeasurementType: req.MeasurementType,
		Notes:           req.Notes,
		MealInfo:        req.MealInfo,
		MedicationInfo:  req.MedicationInfo,
		Indicator:       indicator,
	}

	if err := s.repo.CreateBloodSugarLog(ctx, &log); err != nil {
		return entities.BloodSugarLog{}, err
	}
	return log, nil
}

func (s *healthLogService) GetBloodSugarLogs(ctx context.Context, userID uuid.UUID) ([]entities.BloodSugarLog, error) {
	return s.repo.GetBloodSugarLogs(ctx, userID)
}

func (s *healthLogService) DeleteBloodSugar(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteBloodSugarLog(ctx, id, userID)
}

func (s *healthLogService) CreateBloodPressure(ctx context.Context, userID uuid.UUID, req dto.CreateBloodPressureRequest) (entities.BloodPressureLog, error) {
	classification := classifyBP(req.Systolic, req.Diastolic)

	log := entities.BloodPressureLog{
		UserID:         userID,
		RecordedAt:     req.RecordedAt,
		Systolic:       req.Systolic,
		Diastolic:      req.Diastolic,
		Pulse:          req.Pulse,
		Posture:        req.Posture,
		Classification: classification,
		Notes:          req.Notes,
	}

	if err := s.repo.CreateBloodPressureLog(ctx, &log); err != nil {
		return entities.BloodPressureLog{}, err
	}
	return log, nil
}

func (s *healthLogService) GetBloodPressureLogs(ctx context.Context, userID uuid.UUID) ([]entities.BloodPressureLog, error) {
	return s.repo.GetBloodPressureLogs(ctx, userID)
}

func (s *healthLogService) DeleteBloodPressure(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteBloodPressureLog(ctx, id, userID)
}

func (s *healthLogService) CreateWeight(ctx context.Context, userID uuid.UUID, req dto.CreateWeightRequest) (entities.WeightLog, error) {
	var bmi float64
	var classification string
	if req.HeightCM > 0 {
		heightM := req.HeightCM / 100
		bmi = req.WeightKG / (heightM * heightM)
		bmi = math.Round(bmi*10) / 10
		classification = classifyBMI(bmi)
	}

	log := entities.WeightLog{
		UserID:            userID,
		RecordedAt:        req.RecordedAt,
		WeightKG:          req.WeightKG,
		HeightCM:          req.HeightCM,
		BMI:               bmi,
		BMIClassification: classification,
		BodyFatPercentage: req.BodyFatPercentage,
		Notes:             req.Notes,
	}

	if err := s.repo.CreateWeightLog(ctx, &log); err != nil {
		return entities.WeightLog{}, err
	}
	return log, nil
}

func (s *healthLogService) GetWeightLogs(ctx context.Context, userID uuid.UUID) ([]entities.WeightLog, error) {
	return s.repo.GetWeightLogs(ctx, userID)
}

func (s *healthLogService) DeleteWeight(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteWeightLog(ctx, id, userID)
}

func classifyBP(sys, dia int) string {
	if sys >= 180 || dia >= 120 {
		return "crisis"
	}
	if sys >= 160 || dia >= 100 {
		return "hypertension stage 2"
	}
	if sys >= 140 || dia >= 90 {
		return "hypertension stage 1"
	}
	if sys >= 120 || dia > 80 {
		return "prehypertension"
	}
	if sys < 90 || dia < 60 {
		return "low"
	}
	return "normal"
}

func classifyBMI(bmi float64) string {
	if bmi < 18.5 {
		return "underweight"
	}
	if bmi < 25 {
		return "normal"
	}
	if bmi < 30 {
		return "overweight"
	}
	return "obese"
}
