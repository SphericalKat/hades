package coupon

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Service interface {
	SaveCoupon(*entities.Coupon) error
	RedeemCoupon(uint, string) error
	DeleteCoupon(uint) error
	GetCoupons(uint, uint) ([]entities.Coupon, error)
	VerifyCoupon(uint, uint) (bool, error)
}

type couponSvc struct {
	repo Repository
}

func NewCouponService(r Repository) Service {
	return &couponSvc{repo: r}
}

func (s *couponSvc) SaveCoupon(coupons *entities.Coupon) error {
	return s.repo.SaveCoupon(coupons)
}

func (s *couponSvc) RedeemCoupon(couponId uint, regNo string) error {
	return s.repo.RedeemCoupon(couponId, regNo)
}

func (s *couponSvc) DeleteCoupon(couponId uint) error {
	return s.repo.DeleteCoupon(couponId)
}

func (s *couponSvc) GetCoupons(eventId uint, day uint) ([]entities.Coupon, error) {
	return s.repo.GetCoupons(eventId, day)
}

func (s *couponSvc) VerifyCoupon(eventId uint, couponId uint) (bool, error) {
	return s.repo.VerifyCoupon(eventId, couponId)
}
