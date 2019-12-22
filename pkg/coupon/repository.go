package coupon

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	SaveCoupon(*entities.Coupon) error
	RemoveCouponParticipant(uint, string) error
	DeleteCoupon(uint) error
	GetCoupons(uint) ([]entities.Coupon, error)
	AddCouponsToAll(uint) error
	VerifyCoupon(uint, uint) (bool, error)
}
