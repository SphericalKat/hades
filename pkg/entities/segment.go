package entities

import "time"

type EventSegment struct {
	EventID             uint          `json:"event_id" gorm:"primary_key"`
	Day                 uint          `json:"day" gorm:"primary_key"`
	PresentParticipants []Participant `json:"-" gorm:"many2many:participant_event_segment"`
	DeletedAt           *time.Time    `json:"-" sql:"index"`
}
