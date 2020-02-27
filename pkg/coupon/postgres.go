package coupon

import (
	"log"

	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}

func (r *repo) SaveCoupon(coupon *entities.Coupon) error {
	tx := r.DB.Begin()
	eve := &entities.Event{ID: coupon.EventId}

	switch tx.Where("id = ?", coupon.EventId).Find(eve).Error {
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return pkg.ErrNotFound
	case nil:
		if coupon.Day > eve.Days {
			return pkg.ErrInvalidSlug
		}
		if err := tx.Save(coupon).Error; err != nil {
			tx.Rollback()
			return pkg.ErrDatabase
		}
		tx.Commit()
		return nil
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) DeleteCoupon(couponId uint) error {
	err := r.DB.Delete(&entities.Coupon{CouponId: couponId}).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	case nil:
		return nil
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) RedeemCoupon(couponId uint, email string) error {
	tx := r.DB.Begin()
	p := &entities.Participant{Email: email}
	c := &entities.Coupon{}
	var parts []entities.Participant

	// Find the coupon in db
	if err := tx.Where("coupon_id = ?", couponId).Find(c).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	// Find participant in db
	if err := tx.Where("email = ?", email).Find(p).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	// Check if already redeemed
	if err := tx.Model(c).Association("Participants").Find(&parts).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	if parts == nil {
		return pkg.ErrNotFound
	}
	for _, part := range parts {
		if part.Email == email {
			return pkg.ErrAlreadyExists
		}
	}

	// verify that the participant is part of the event
	if err := tx.Model(p).Association("Events").Find(&(p.Events)).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			log.Println(err)
			return pkg.ErrDatabase
		}
	}
	flag := false
	for _, event := range p.Events {
		if event.ID == c.EventId {
			flag = true
		}
	}

	if !flag {
		return pkg.ErrInvalidSlug
	}

	// Redeem coupon by adding to association
	if err := tx.Model(c).Association("Participants").Append(p).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			log.Println(err)
			return pkg.ErrDatabase
		}
	}

	tx.Commit()
	return nil
}

func (r *repo) GetCoupons(eventId uint, day uint) ([]entities.Coupon, error) {
	var coupons []entities.Coupon
	err := r.DB.Model(entities.Coupon{}).Where("event_id = ?", eventId).Where("day = ?", day).Find(&coupons).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	case nil:
		return coupons, nil
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) VerifyCoupon(eventId uint, couponId uint) (bool, error) {
	coup := &entities.Coupon{CouponId: couponId}
	err := r.DB.Where("coupon_id = ?", couponId).Find(coup).Error

	if err == nil {
		if coup.EventId != eventId {
			return false, nil
		}
		return true, nil
	} else if err == gorm.ErrRecordNotFound {
		return false, pkg.ErrNotFound
	} else {
		return false, pkg.ErrDatabase
	}
}
