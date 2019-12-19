package guest

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/ATechnoHazard/hades-2/pkg/entities"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB:db}
}

func (r *repo) FindGuestEvent(email string, eventId uint) (*entities.Guest, error) {
	eve := &entities.Event{}

	err := r.DB.Where("event_id = ?", eventId).Find(eve).Error
	switch err{
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	case nil:
		for _, gus := range eve.Guests {
			if gus.Email == email {
				return &gus, nil
			}
		}
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) FindAllGuestEvent(eventId uint) ([]entities.Guest, error) {
	eve := &entities.Event{}
	err := r.DB.Where("event_id = ?", eventId).Find(eve).Error
	switch err {
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	case nil:
		return eve.Guests, nil
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) SaveGuestEvent(guest *entities.Guest, eventId uint) error {
	tx := r.DB.Begin()
	e := &entities.Event{ID:eventId}

	if tx.Find(e).Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	}

	err := tx.Unscoped().Save(guest).Error
	if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}

	err = tx.Model(guest).Association("Events").Append(e).Error
	if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}

	tx.Commit()
	return nil
}

func (r *repo) Delete(email string) error {
	err := r.DB.Where("email = ?", email).Delete(&entities.Guest{}).Error
	if err != nil {
		return pkg.ErrDatabase
	}
	return nil
}

func (r *repo) RemoveGuestEvent(emailId string, eventID uint) error {
	tx := r.DB.Begin()
	e := &entities.Event{ID: eventID}
	p := &entities.Guest{Email:emailId}

	if tx.Find(e).Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	}
	if tx.Find(p).Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	}

	err := tx.Model(p).Association("Events").Delete(e).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	tx.Commit()
	return nil
}
