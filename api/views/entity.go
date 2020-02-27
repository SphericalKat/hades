package views

import (
	"github.com/ATechnoHazard/hades-2/pkg/entities"
)

type Participant struct {
	Name        string `json:"name"`
	RegNo       string `json:"reg_no"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
	EventId     uint   `json:"event_id"`
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

type Guest struct {
	Name           string `json:"name"`
	Email          string `json:"email" gorm:"primary_key;"`
	PhoneNumber    string `json:"phone_number"`
	Gender         string `json:"gender"`
	Stake          string `json:"stake"`
	LocationOfStay string `json:"location_of_stay"`
	EventId        uint   `json:"event_id"`
}

func (g *Guest) Transform() *entities.Guest {
	return &entities.Guest{
		Name:           g.Name,
		Email:          g.Email,
		PhoneNumber:    g.PhoneNumber,
		Gender:         g.Gender,
		Stake:          g.Stake,
		LocationOfStay: g.LocationOfStay,
		Events:         nil,
		DeletedAt:      nil,
	}
}

type CouponParticipantComposite struct {
	CouponID uint   `json:"coupon_id"`
	RegNo    string `json:"reg_no"`
	EventID  uint   `json:"event_id"`
	Email    string `json:"email"`
}

type SegmentParticipantComposite struct {
	Day     uint   `json:"day"`
	RegNo   string `json:"reg_no"`
	EventID uint   `json:"event_id"`
	Email   string `json:"email"`
}
