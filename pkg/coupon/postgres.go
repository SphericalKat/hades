package coupon

import (
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

	switch tx.Find(eve).Error {
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return pkg.ErrNotFound
	case nil:
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

func (r *repo) RemoveCouponParticipant(couponId uint, regNo string) error {
	tx := r.DB.Begin()
	p := &entities.Participant{RegNo: regNo}

	switch tx.Find(p).Error {
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound

	case nil:
		err := tx.Model(p).Association("Coupons").Delete(&entities.Coupon{CouponId: couponId}).Error

		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		case nil:
			return nil
		default:
			return pkg.ErrDatabase
		}

	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) GetCoupons(eventId uint) ([]entities.Coupon, error) {
	var coupons []entities.Coupon
	err := r.DB.Model(entities.Coupon{}).Where("event_id = ?", eventId).Find(&coupons).Error

	switch err {
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	case nil:
		return coupons, nil
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) AddCouponsToAll(eventId uint) error {
	eve := &entities.Event{ID: eventId}
	var peeps []entities.Participant
	var coups []entities.Coupon

	tx := r.DB.Begin()
	err := tx.Find(eve).Error

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	} else if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}

	// populate coupons by event id.
	err = tx.Model(entities.Coupon{}).Where("event_id = ?", eventId).Find(&coups).Error
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	} else if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}

	// populate participants in the event.
	err = tx.Model(eve).Association("Attendees").Find(&peeps).Error

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	} else if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}

	for _, coup := range coups {
		for _, peep := range peeps {
			err = tx.Model(peep).Association("Coupons").Append(coup).Error
			if err == gorm.ErrRecordNotFound {
				tx.Rollback()
				return pkg.ErrNotFound
			} else if err != nil {
				tx.Rollback()
				return pkg.ErrDatabase
			}
		}
	}

	return nil
}
