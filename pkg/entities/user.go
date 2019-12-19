package entities

import (
	"time"
)

type User struct {
	FirstName     string                      `json:"first_name"`
	LastName      string                      `json:"last_name"`
	Password      string                      `json:"password"`
	Email         string                      `json:"email" gorm:"primary_key"`
	PhoneNumber   string                      `json:"phone_number"`
	Linkedin      string                      `json:"linkedin"`
	Facebook      string                      `json:"facebook"`
	Description   string                      `json:"description"`
	CreatedAt     time.Time                   `json:"created_at"`
	DeviceToken   string                      `json:"device_token"`
	Organizations []Organization `json:"-" gorm:"many2many:user_orgs;"`
}
