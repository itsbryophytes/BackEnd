package entities

import (
	"time"

	"github.com/google/uuid"
)

type BloodSugarLog struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RecordedAt      time.Time `gorm:"type:timestamptz;not null" json:"recorded_at"`
	GlucoseValue    float64   `gorm:"not null" json:"glucose_value"`
	MeasurementType string    `gorm:"not null" json:"measurement_type"`
	Notes           string    `json:"notes"`
	MealInfo        string    `json:"meal_info"`
	MedicationInfo  string    `json:"medication_info"`
	Indicator       string    `json:"indicator"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type BloodPressureLog struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RecordedAt     time.Time `gorm:"type:timestamptz;not null" json:"recorded_at"`
	Systolic       int       `gorm:"not null" json:"systolic"`
	Diastolic      int       `gorm:"not null" json:"diastolic"`
	Pulse          int       `json:"pulse"`
	Posture        string    `json:"posture"`
	Classification string    `json:"classification"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type WeightLog struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RecordedAt        time.Time `gorm:"type:timestamptz;not null" json:"recorded_at"`
	WeightKG          float64   `gorm:"not null" json:"weight_kg"`
	HeightCM          float64   `json:"height_cm"`
	BMI               float64   `json:"bmi"`
	BMIClassification string    `json:"bmi_classification"`
	BodyFatPercentage float64   `json:"body_fat_percentage"`
	Notes             string    `json:"notes"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
}