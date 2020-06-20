package entities

import (
	"time"
)

type Participant struct {
	Name          string         `json:"name"`
	RegNo         string         `json:"reg_no"`
	Email         string         `json:"email" gorm:"primary_key"`
	PhoneNumber   string         `json:"phone_number"`
	Gender        string         `json:"gender"`
	DeletedAt     *time.Time     `json:"-" sql:"index"`
	Events        []Event        `json:"-" gorm:"many2many:participant_events;"`
	Coupons       []Coupon       `json:"-" gorm:"many2many:coupon_participant;"`
	EventSegments []EventSegment `json:"-" gorm:"many2many:participant_event_segment"`
}
