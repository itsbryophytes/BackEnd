package dto

import "time"

type CreateBloodSugarRequest struct {
	RecordedAt      time.Time `json:"recorded_at" binding:"required"`
	GlucoseValue    float64   `json:"glucose_value" binding:"required"`
	MeasurementType string    `json:"measurement_type" binding:"required"`
	Notes           string    `json:"notes"`
	MealInfo        string    `json:"meal_info"`
	MedicationInfo  string    `json:"medication_info"`
}

type CreateBloodPressureRequest struct {
	RecordedAt time.Time `json:"recorded_at" binding:"required"`
	Systolic   int       `json:"systolic" binding:"required"`
	Diastolic  int       `json:"diastolic" binding:"required"`
	Pulse      int       `json:"pulse"`
	Posture    string    `json:"posture"`
	Notes      string    `json:"notes"`
}

type CreateWeightRequest struct {
	RecordedAt        time.Time `json:"recorded_at" binding:"required"`
	WeightKG          float64   `json:"weight_kg" binding:"required"`
	HeightCM          float64   `json:"height_cm"`
	BodyFatPercentage float64   `json:"body_fat_percentage"`
	Notes             string    `json:"notes"`
}
