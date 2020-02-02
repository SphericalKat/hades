package coupon

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Repository interface {
	SaveCoupon(*entities.Coupon) error
	RedeemCoupon(uint, string) error
	DeleteCoupon(uint) error
	GetCoupons(uint, uint) ([]entities.Coupon, error)
	VerifyCoupon(uint, uint) (bool, error)
}
