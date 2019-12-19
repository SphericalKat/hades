package participant

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

func (r *repo) FindAll() ([]entities.Participant, error) {
	var participants []entities.Participant
	err := r.DB.Model(entities.Participant{}).Find(&participants).Error
	switch err {
	case nil:
		return participants, nil
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) FindByRegNo(regNo string, eventID uint) (*entities.Participant, error) {
	e := &entities.Event{ID: eventID}

	err := r.DB.Find(e).Association("Attendees").Find(&e.Attendees).Error
	switch err {
	case nil:
		for _, p := range e.Attendees {
			if p.RegNo == regNo {
				return &p, nil
			}
		}
		return nil, pkg.ErrNotFound
	case gorm.ErrRecordNotFound:
		return nil, pkg.ErrNotFound
	default:
		return nil, pkg.ErrDatabase
	}
}

func (r *repo) Save(participant *entities.Participant, eventID uint) error {
	tx := r.DB.Begin()
	e := &entities.Event{ID: eventID}

	if tx.Find(e).Error == gorm.ErrRecordNotFound {
		tx.Rollback()
		return pkg.ErrNotFound
	}

	err := tx.Unscoped().Save(participant).Error
	if err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	err = tx.Model(participant).Association("Events").Append(e).Error
	if err != nil {
		tx.Rollback()
		return pkg.ErrDatabase
	}
	tx.Commit()
	return nil
}

func (r *repo) Delete(regNo string) error {
	err := r.DB.Where("reg_no = ?", regNo).Delete(&entities.Participant{}).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return pkg.ErrNotFound
	default:
		return pkg.ErrDatabase
	}
}

func (r *repo) RemoveAttendeeEvent(regNo string, eventID uint) error {
	tx := r.DB.Begin()
	e := &entities.Event{ID: eventID}
	p := &entities.Participant{RegNo: regNo}

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
