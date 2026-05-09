package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserHealthProfile struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID             uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	DateOfBirth        *time.Time `gorm:"type:date" json:"date_of_birth"`
	BiologicalSex      string    `gorm:"type:varchar(10)" json:"biological_sex"`
	HeightCm           *float64  `gorm:"type:float" json:"height_cm"`
	WeightKg           *float64  `gorm:"type:float" json:"weight_kg"`
	BloodType          string    `gorm:"type:varchar(5)" json:"blood_type"`
	SmokingStatus      string    `gorm:"type:varchar(20)" json:"smoking_status"`
	ExistingConditions string    `gorm:"type:text" json:"existing_conditions"`
	CurrentMedications string    `gorm:"type:text" json:"current_medications"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
