package participant

import "github.com/ATechnoHazard/hades-2/pkg/event"

type Participant struct {
	Name        string        `json:"name"`
	RegNo       string        `json:"reg_no" gorm:"primary_key"`
	Email       string        `json:"email"`
	PhoneNumber string        `json:"phone_number"`
	Gender      string        `json:"gender"`
	Events      []event.Event `json:"event_name" gorm:"many2many:participant_events;"`
}
