package event

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

func (r *repo) Find(eventID uint) (*entities.Event, error) {
	event := &entities.Event{}
	tx := r.DB.Begin()
	err := tx.Where("id = ?", eventID).Find(event).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			tx.Rollback()
			return nil, pkg.ErrNotFound
		default:
			tx.Rollback()
			return nil, pkg.ErrDatabase
		}
	}

	err = tx.Find(event).Association("Attendees").Find(&event.Attendees).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, pkg.ErrNotFound
		default:
			return nil, pkg.ErrDatabase
		}
	}

	err = tx.Find(event).Association("Guests").Find(&event.Guests).Error
	switch err {
	case nil:
		tx.Commit()
		return event, nil
	case gorm.ErrRecordNotFound:
		tx.Rollback()
		return nil, pkg.ErrNotFound
	default:
		tx.Rollback()
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) Save(event *entities.Event) (*entities.Event, error) {
	if event.Days == 0 {
		return nil, pkg.ErrInvalidSlug
	}
	tx := r.DB.Begin()
	err := tx.Save(event).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, pkg.ErrNotFound
		default:
			return nil, pkg.ErrDatabase
		}
	}

	for i := uint(1); i <= event.Days; i++ {
		s := &entities.EventSegment{EventID: event.ID, Day: i}
		err := tx.Create(s).Error
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				tx.Rollback()
				return nil, pkg.ErrNotFound
			default:
				tx.Rollback()
				return nil, pkg.ErrDatabase
			}
		}
	}
	tx.Commit()
	return event, err
}

func (r *repo) Delete(eventID uint) error {
	err := r.DB.Where("id = ?", eventID).Delete(&entities.Event{}).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}
