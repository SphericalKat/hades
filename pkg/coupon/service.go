package coupon

import "github.com/ATechnoHazard/hades-2/pkg/entities"

type Service interface {
	SaveCoupon(*entities.Coupon) error
	RemoveCouponParticipant(uint, string) error
	DeleteCoupon(uint) error
	GetCoupons(uint) ([]entities.Coupon, error)
	AddCouponsToAll(uint) error
	VerifyCoupon(uint, uint) (bool, error)
}

type couponSvc struct {
	repo Repository
}

func NewCouponService(r Repository) Service {
	return &couponSvc{repo:r}
}

func (s *couponSvc) SaveCoupon(coupons *entities.Coupon) error {
	return s.repo.SaveCoupon(coupons)
}

func (s *couponSvc) RemoveCouponParticipant(couponId uint, regNo string) error {
	return s.repo.RemoveCouponParticipant(couponId, regNo)
}

func (s *couponSvc) DeleteCoupon(couponId uint) error {
	return s.repo.DeleteCoupon(couponId)
}

func (s *couponSvc) GetCoupons(eventId uint) ([]entities.Coupon, error) {
	return s.repo.GetCoupons(eventId)
}

func (s *couponSvc) AddCouponsToAll(eventId uint) error {
	return s.repo.AddCouponsToAll(eventId)
}

func (s *couponSvc) VerifyCoupon(eventId uint, couponId uint) (bool, error) {
	return s.repo.VerifyCoupon(eventId, couponId)
}