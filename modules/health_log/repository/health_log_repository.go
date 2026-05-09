package repository

import (
	"context"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HealthLogRepository interface {
	CreateBloodSugarLog(ctx context.Context, log *entities.BloodSugarLog) error
	GetBloodSugarLogs(ctx context.Context, userID uuid.UUID) ([]entities.BloodSugarLog, error)
	DeleteBloodSugarLog(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	CreateBloodPressureLog(ctx context.Context, log *entities.BloodPressureLog) error
	GetBloodPressureLogs(ctx context.Context, userID uuid.UUID) ([]entities.BloodPressureLog, error)
	DeleteBloodPressureLog(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

	CreateWeightLog(ctx context.Context, log *entities.WeightLog) error
	GetWeightLogs(ctx context.Context, userID uuid.UUID) ([]entities.WeightLog, error)
	DeleteWeightLog(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type healthLogRepository struct {
	db *gorm.DB
}

func NewHealthLogRepository(db *gorm.DB) HealthLogRepository {
	return &healthLogRepository{db: db}
}

func (r *healthLogRepository) CreateBloodSugarLog(ctx context.Context, log *entities.BloodSugarLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *healthLogRepository) GetBloodSugarLogs(ctx context.Context, userID uuid.UUID) ([]entities.BloodSugarLog, error) {
	var logs []entities.BloodSugarLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("recorded_at DESC").Find(&logs).Error
	return logs, err
}

func (r *healthLogRepository) DeleteBloodSugarLog(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&entities.BloodSugarLog{}).Error
}

func (r *healthLogRepository) CreateBloodPressureLog(ctx context.Context, log *entities.BloodPressureLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *healthLogRepository) GetBloodPressureLogs(ctx context.Context, userID uuid.UUID) ([]entities.BloodPressureLog, error) {
	var logs []entities.BloodPressureLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("recorded_at DESC").Find(&logs).Error
	return logs, err
}

func (r *healthLogRepository) DeleteBloodPressureLog(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&entities.BloodPressureLog{}).Error
}

func (r *healthLogRepository) CreateWeightLog(ctx context.Context, log *entities.WeightLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *healthLogRepository) GetWeightLogs(ctx context.Context, userID uuid.UUID) ([]entities.WeightLog, error) {
	var logs []entities.WeightLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("recorded_at DESC").Find(&logs).Error
	return logs, err
}

func (r *healthLogRepository) DeleteWeightLog(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&entities.WeightLog{}).Error
}
