package model

import (
	"time"

	"github.com/dioncodes/go-waitlist/internal/db"
	"github.com/dioncodes/go-waitlist/internal/errors"
	"gorm.io/datatypes"
)

type Registration struct {
	ID                    int             `json:"-" gorm:"column:id;primaryKey;autoIncrement;type:int unsigned"`
	Email                 string          `json:"email" gorm:"column:email;uniqueIndex;type:varchar(100);not null"`
	RegistrationDate      time.Time       `json:"registrationDate" gorm:"column:registration_date;type:datetime;autoCreateTime;default:NOW()"`
	ConfirmationToken     string          `json:"confirmationToken" gorm:"column:confirmation_token;type:varchar(32)"`
	Confirmed             bool            `json:"confirmed" gorm:"column:confirmed;not null;default:0"`
	AdditionalInformation datatypes.JSON  `json:"additionalInformation,omitempty" gorm:"column:additional_information;type:json"`
}

func (registration *Registration) Save() error {
	if registration.ConfirmationToken == "" {
		registration.ConfirmationToken = GenerateToken(16)
	}

	result := db.Conn().Save(registration)

	return result.Error
}

func (registration *Registration) Delete() error {
	result := db.Conn().Delete(registration)

	return result.Error
}

func GetRegistrationByEmail(email string) (Registration, error) {
	var registration Registration

	db.Conn().Where("email = ?", email).First(&registration)

	if registration.ID == 0 {
		return registration, errors.NotFound()
	}

	return registration, nil
}

func GetRegistrationById(id string) (Registration, error) {
	var registration Registration

	db.Conn().Where("id = ?", id).First(&registration)

	if registration.ID == 0 {
		return registration, errors.NotFound()
	}

	return registration, nil
}

func GetAllRegistrations() ([]Registration, error) {
	var registrations []Registration

	result := db.Conn().Find(&registrations)

	return registrations, result.Error
}
