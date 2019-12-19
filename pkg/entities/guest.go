package entities

import "time"

type Guest struct {
	Name           string     `json:"name"`
	Email          string     `json:"email" gorm:"primary_key;"`
	PhoneNumber    string     `json:"phone_number"`
	Gender         string     `json:"gender"`
	Stake          string     `json:"stake"`
	LocationOfStay string     `json:"location_of_stay"`
	Events         []Event    `json:"-" gorm:"many2many:guest_event;"`
	DeletedAt      *time.Time `json:"-" sql:"index"`
}
