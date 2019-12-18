package participant

import (
	"github.com/ATechnoHazard/hades-2/pkg/event"
	"time"
)

type Participant struct {
	Name        string        `json:"name"`
	RegNo       string        `json:"reg_no" gorm:"primary_key"`
	Email       string        `json:"email"`
	PhoneNumber string        `json:"phone_number"`
	Gender      string        `json:"gender"`
	DeletedAt   *time.Time    `json:"-" sql:"index"`
	Events      []event.Event `json:"-" gorm:"many2many:participant_events;"`
	EventId     uint          `json:"event_id" gorm:"-"`
}
