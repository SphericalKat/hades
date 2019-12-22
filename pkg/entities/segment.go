package entities

import "time"

type EventSegment struct {
	EventID             uint          `json:"event_id"`
	SegmentID           uint          `json:"segment_id" gorm:"primary_key;AUTO_INCREMENT"`
	Name                string        `json:"name"`
	PresentParticipants []Participant `json:"-" gorm:"many2many:participant_event_segment"`
	DeletedAt           *time.Time    `json:"-" sql:"index"`
}

type ESParticipant struct {
	Participant Participant `json:"participant"`
	IsPresent   bool        `json:"is_present"`
}
