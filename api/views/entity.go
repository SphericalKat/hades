package views

import (
	"github.com/ATechnoHazard/hades-2/pkg/entities"
)

type Participant struct {
	Name        string        `json:"name"`
	RegNo       string        `json:"reg_no"`
	Email       string        `json:"email"`
	PhoneNumber string        `json:"phone_number"`
	Gender      string        `json:"gender"`
	EventId     uint          `json:"event_id"`
}

func (p Participant) Transform() *entities.Participant {
	return &entities.Participant{
		Name:        p.Name,
		RegNo:       p.RegNo,
		Email:       p.Email,
		PhoneNumber: p.PhoneNumber,
		Gender:      p.Gender,
		DeletedAt:   nil,
		Events:      nil,
	}
}
