package user

import "github.com/ATechnoHazard/hades-2/pkg/organization"

type User struct {
	FirstName     string                      `json:"first_name"`
	LastName      string                      `json:"last_name"`
	Password      string                      `json:"password"`
	Email         string                      `json:"email" gorm:"primary_key"`
	PhoneNumber   string                      `json:"phone_number"`
	Linkedin      string                      `json:"linkedin"`
	Facebook      string                      `json:"facebook"`
	Description   string                      `json:"description"`
	CreatedAt     string                      `json:"created_at"`
	DeviceToken   string                      `json:"device_token"`
	Organizations []organization.Organization `json:"-" gorm:"many2many:user_orgs;"`
}
