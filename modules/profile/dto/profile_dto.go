package dto

import "errors"

const (
	MESSAGE_SUCCESS_GET_PROFILE    = "success get health profile"
	MESSAGE_SUCCESS_UPDATE_PROFILE = "success update health profile"
	MESSAGE_FAILED_GET_PROFILE     = "failed get health profile"
	MESSAGE_FAILED_UPDATE_PROFILE  = "failed update health profile"
)

var (
	ErrProfileNotFound = errors.New("health profile not found")
)

type HealthProfileRequest struct {
	DateOfBirth        string   `json:"date_of_birth"`
	BiologicalSex      string   `json:"biological_sex"`
	HeightCm           *float64 `json:"height_cm"`
	WeightKg           *float64 `json:"weight_kg"`
	BloodType          string   `json:"blood_type"`
	SmokingStatus      string   `json:"smoking_status"`
	ExistingConditions string   `json:"existing_conditions"`
	CurrentMedications string   `json:"current_medications"`
}

type HealthProfileResponse struct {
	ID                 string   `json:"id"`
	UserID             string   `json:"user_id"`
	DateOfBirth        string   `json:"date_of_birth"`
	BiologicalSex      string   `json:"biological_sex"`
	HeightCm           *float64 `json:"height_cm"`
	WeightKg           *float64 `json:"weight_kg"`
	BloodType          string   `json:"blood_type"`
	SmokingStatus      string   `json:"smoking_status"`
	ExistingConditions string   `json:"existing_conditions"`
	CurrentMedications string   `json:"current_medications"`
}
