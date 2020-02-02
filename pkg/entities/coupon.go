package entities

import "time"

type Coupon struct {
	CouponId     uint          `json:"coupon_id" gorm:"primary_key;AUTO_INCREMENT"`
	EventId      uint          `json:"event_id"`
	Day          uint          `json:"day"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	DeletedAt    *time.Time    `json:"-" sql:"index"`
	Participants []Participant `json:"-" gorm:"many2many:coupon_participant"`
}
