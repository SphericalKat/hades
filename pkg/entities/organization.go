package entities

import (
	"time"
)

type Organization struct {
	ID           uint          `json:"org_id" gorm:"primary_key;AUTO_INCREMENT"`
	Name         string        `json:"name"`
	Location     string        `json:"location"`
	Description  string        `json:"description"`
	Tag          string        `json:"tag"`
	Website      string        `json:"website"`
	CreatedAt    time.Time     `json:"created_at"`
	DeletedAt    *time.Time    `sql:"index"`
	Events       []Event       `json:"-"`
	JoinRequests []JoinRequest `json:"-"`
	Users        []User        `json:"-" gorm:"many2many:user_orgs;"`
}

type JoinRequest struct {
	OrganizationID uint   `gorm:"primary_key" json:"org_id"`
	Email          string `gorm:"primary_key" json:"email"`
}
