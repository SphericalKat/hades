package event

import (
	"github.com/ATechnoHazard/hades-2/pkg"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}


func (r *repo) Find(eventID uint) (*Event, error) {
	event := &Event{ID: eventID}
	err := r.DB.Find(event).Error
	switch err {
	case nil:
		return event, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) Save(event *Event) error {
	err := r.DB.Save(event).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) Delete(eventID uint) error {
	err := r.DB.Where("event_id = ?", eventID).Delete(&Event{}).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}
