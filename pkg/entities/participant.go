package entities

import (
	"time"
)

type Participant struct {
	Name        string     `json:"name"`
	RegNo       string     `json:"reg_no" gorm:"primary_key"`
	Email       string     `json:"email"`
	PhoneNumber string     `json:"phone_number"`
	Gender      string     `json:"gender"`
	DeletedAt   *time.Time `json:"-" sql:"index"`
	Events      []Event    `json:"-" gorm:"many2many:participant_events;"`
}
